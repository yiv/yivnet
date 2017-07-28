package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	ketcd "github.com/go-kit/kit/sd/etcd"
	"github.com/go-kit/kit/tracing/opentracing"

	"github.com/yiv/yivgame/game/gamer"
	"github.com/yiv/yivgame/game/pb"
	"github.com/yiv/yivgame/game/service"
)

var (
	grpcAddr     = flag.String("grpc.addr", ":10020", "gRPC (HTTP) listen address")
	debugAddr    = flag.String("debug.addr", ":10021", "Debug and metrics listen address")
	zipkinAddr   = flag.String("zipkin.addr", "", "Enable Zipkin tracing via a Zipkin HTTP Collector endpoint")
	etcdAddr     = flag.String("etcd.addr", "http://etcd.yivgame.com:2379", "Consul agent address")
	serviceName  = flag.String("service.name", "gameserver", "the name of this service in service discovery")
	userCenter   = flag.String("userCenter.name", "usercenter", "the name of this service in service discovery")
	retryMax     = flag.Int("retry.max", 3, "per-request retries to different instances")
	retryTimeout = flag.Duration("retry.timeout", 500*time.Millisecond, "per-request timeout, including retries")
	//first class room config
	frBootbet = flag.Int64("fr.bb", 100, "boot bet coins of first class room")
)

func main() {
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	logger.Log("msg", "hello")
	defer logger.Log("msg", "goodbye")

	// Mechanical domain.
	errc := make(chan error)
	// Interrupt handler.
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// Service discovery domain
	ctx := context.Background()
	etcdClient, err := ketcd.NewClient(ctx, []string{*etcdAddr}, ketcd.ClientOptions{})
	if err != nil {
		panic(err)
	}
	registrar := ketcd.NewRegistrar(etcdClient, ketcd.Service{Key: *serviceName, Value: *grpcAddr}, logger)
	registrar.Register()
	defer registrar.Deregister()

	// Metrics domain.
	var (
		requestCount   metrics.Counter
		requestLatency metrics.Histogram
		fieldKeys      []string
	)
	{
		// Business level metrics.
		fieldKeys = []string{"method", "error"}
		requestCount = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "gameserver",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys)
		requestLatency = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "usercenter",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys)
	}
	var duration metrics.Histogram
	{
		// Transport level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "gameserver",
			Name:      "request_duration_ns",
			Help:      "Request duration in nanoseconds.",
		}, []string{"method", "success"})
	}

	// Tracing domain.
	var tracer stdopentracing.Tracer
	{

		logger := log.With(logger, "tracer", "ZipkinHTTP")
		logger.Log("addr", *zipkinAddr)

		// endpoint typically looks like: http://zipkinhost:9411/api/v1/spans
		collector, err := zipkin.NewHTTPCollector(*zipkinAddr)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		defer collector.Close()

		tracer, err = zipkin.NewTracer(
			zipkin.NewRecorder(collector, false, "localhost:10003", "addsvc"),
		)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}

	}
	// Business domain.
	uc, err := service.NewUserCenter(*userCenter, []string{*etcdAddr}, *retryMax, *retryTimeout, logger)
	if err != nil {
		errc <- err
		return
	}
	firRoomOp := gamer.RoomOptions{
		RoomClass: gamer.FirstClassRoom,
		BootBet:   *frBootbet,
	}
	var gameService service.GameService
	{
		logger := log.With(logger, "service", "gameService")
		gameService = service.NewGameService(firRoomOp, uc, logger)
		gameService = service.ServiceLoggingMiddleware(logger)(gameService)
		gameService = service.ServiceInstrumentingMiddleware(requestCount, requestLatency)(gameService)
	}
	// Endpoint domain.
	var sendChatEndpoint endpoint.Endpoint
	{
		method := "sendChat"
		duration := duration.With("method", method)
		logger := log.With(logger, "method", method)
		sendChatEndpoint = service.MakeSendChatEndpoint(gameService)
		sendChatEndpoint = opentracing.TraceServer(tracer, method)(sendChatEndpoint)
		sendChatEndpoint = service.EndpointInstrumentingMiddleware(duration)(sendChatEndpoint)
		sendChatEndpoint = service.EndpointLoggingMiddleware(logger)(sendChatEndpoint)
	}
	endpoints := service.Endpoints{
		Logger:           logger,
		SendChatEndpoint: sendChatEndpoint,
	}

	// gRPC transport.
	go func() {
		logger := log.With(logger, "transport", "gRPC")

		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}
		grpcServer := grpc.NewServer()
		grpcHandler := service.MakeGRPCHandler(endpoints, tracer, logger)
		pb.RegisterGameServiceServer(grpcServer, grpcHandler)

		logger.Log("addr", *grpcAddr)
		errc <- grpcServer.Serve(ln)
	}()

	// Debug listener.
	go func() {
		logger := log.With(logger, "transport", "debug")

		m := http.NewServeMux()
		m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		m.Handle("/metrics", promhttp.Handler())

		logger.Log("addr", *debugAddr)
		errc <- http.ListenAndServe(*debugAddr, m)
	}()
	logger.Log("terminated", <-errc)
}

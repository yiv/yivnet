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

	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	ketcd "github.com/go-kit/kit/sd/etcd"

	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yiv/yivgame/usercenter/cockroach"
	"github.com/yiv/yivgame/usercenter/kafka"
	"github.com/yiv/yivgame/usercenter/pb"
	"github.com/yiv/yivgame/usercenter/service"
)

var (
	grpcAddr    = flag.String("grpc.addr", ":10000", "gRPC listen address")
	debugAddr   = flag.String("debug.addr", ":10001", "Debug and metrics listen address")
	etcdAddr    = flag.String("etcd.addr", "http://etcd.yivgame.com:2379", "Consul agent address")
	kafkaAddr   = flag.String("kafka.addr", "kafka.yivgame.com:9092", "kafka address")
	zipkinAddr  = flag.String("zipkin.addr", "", "Enable Zipkin tracing via a Zipkin HTTP Collector endpoint")
	serviceName = flag.String("service.name", "usercenter", "the name of this service in service discovery")
)

func main() {
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
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
			Namespace: "usercenter",
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
			Namespace: "usercenter",
			Name:      "request_duration_ns",
			Help:      "Request duration in nanoseconds.",
		}, []string{"method", "success"})
	}

	// Setup repositories
	dblogger := log.With(logger, "db", "cockroach")
	dbRepo, err := cockroach.NewDbUserRepo("postgresql://edwin@db.yivgame.com:26257/tp_user?sslmode=disable", dblogger)
	if err != nil {
		logger.Log("cockroach connect err: ", err)
		os.Exit(1)
	}
	mqlogger := log.With(logger, "mq", "kafka")
	mqRepo, err := kafka.NewMQRepo([]string{*kafkaAddr}, mqlogger)
	if err != nil {
		logger.Log("kafka connect err: ", err)
		os.Exit(1)
	}

	// Business domain.
	var basicService service.Service
	{
		basicService = service.NewBasicService(dbRepo, mqRepo, logger)
		basicService = service.ServiceLoggingMiddleware(logger)(basicService)
		basicService = service.ServiceInstrumentingMiddleware(requestCount, requestLatency)(basicService)
	}

	var getUserInfoEndpoint endpoint.Endpoint
	{
		method := "getUserInfo"
		duration := duration.With("method", method)
		logger := log.With(logger, "method", method)
		getUserInfoEndpoint = service.MakeGetUserInfoEndpoint(basicService)
		getUserInfoEndpoint = opentracing.TraceServer(tracer, method)(getUserInfoEndpoint)
		getUserInfoEndpoint = service.EndpointInstrumentingMiddleware(duration)(getUserInfoEndpoint)
		getUserInfoEndpoint = service.EndpointLoggingMiddleware(logger)(getUserInfoEndpoint)
	}
	endpoints := service.Endpoints{
		GetUserInfoEndpoint: getUserInfoEndpoint,
	}

	// gRPC transport.
	go func() {
		logger := log.With(logger, "transport", "gRPC")

		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}
		grpcHandler := service.MakeGRPCHandler(endpoints, tracer, logger)
		s := grpc.NewServer()
		pb.RegisterUserServer(s, grpcHandler)

		logger.Log("addr", *grpcAddr)
		errc <- s.Serve(ln)
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

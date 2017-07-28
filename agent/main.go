package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/yiv/yivgame/agent/service"
	"github.com/yiv/yivgame/game/pb"
)

var (
	webSocketAddr = flag.String("websocket.addr", ":10050", "game agent webSocket address")
	gameServer    = flag.String("grpc.gameserver", ":10020", "game server gRPC server address")
	debugAddr     = flag.String("debug.addr", ":10051", "Debug and metrics listen address")
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

	// gRPC transport.
	conn, err := grpc.Dial(*gameServer, grpc.WithInsecure())
	if err != nil {
		errc <- err
		return
	}
	defer conn.Close()
	client := pb.NewGameServiceClient(conn)

	// Business domain.
	agentService := service.NewAgentService(client, logger)

	// webSocket transport.
	go func() {
		logger := log.With(logger, "transport", "webSocket")
		m := http.NewServeMux()
		m.HandleFunc("/", agentService.WebSocketServer)
		logger.Log("addr", *webSocketAddr)
		errc <- http.ListenAndServe(*webSocketAddr, m)
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

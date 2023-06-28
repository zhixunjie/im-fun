package main

import (
	logicGrpc "github.com/zhixunjie/im-fun/internal/logic/api/grpc"
	logicHttp "github.com/zhixunjie/im-fun/internal/logic/api/http"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/internal/logic/service"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/perf"
	"github.com/zhixunjie/im-fun/pkg/prometheus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// init pprof
	perf.InitPProf(":6062")
	// init prometheus
	prometheus.InitPrometheus(":7061")
	// init config
	if err := conf.InitConfig("cmd/logic/logic.yaml"); err != nil {
		panic(err)
	}
	// init service
	svc := service.New(conf.Conf)
	// init HTTP server
	httpSrv := logicHttp.New(conf.Conf, svc)
	// init GRPC server
	rpcSrv := logicGrpc.New(conf.Conf.RPC.Server, svc)

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logging.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			httpSrv.Close()
			rpcSrv.GracefulStop()
			return
		case syscall.SIGHUP:
			return
		default:
			return
		}
	}
}

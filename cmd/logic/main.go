package main

import (
	"context"
	"github.com/zhixunjie/im-fun/cmd/logic/wire"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/perf"
	"github.com/zhixunjie/im-fun/pkg/prometheus/register"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// init pprof
	perf.InitPProf(":6062")
	// init prometheus
	register.InitProm(":7062")
	// init config
	if err := conf.InitConfig("cmd/logic/logic.yaml"); err != nil {
		panic(err)
	}

	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// init HTTP server
	httpSrv := wire.InitHttp(conf.Conf)
	// init GRPC server
	rpcSrv, deRegister, err := wire.InitGrpc(ctx, conf.Conf)
	if err != nil {
		panic(err)
	}

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logging.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			deRegister()
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

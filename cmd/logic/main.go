package main

import (
	"github.com/sirupsen/logrus"
	logicGrpc "github.com/zhixunjie/im-fun/internal/logic/api/grpc"
	logicHttp "github.com/zhixunjie/im-fun/internal/logic/api/http"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/internal/logic/service"
	"github.com/zhixunjie/im-fun/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.InitLogConfig()
	if err := conf.InitConfig(); err != nil {
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
		logrus.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			httpSrv.Close()
			rpcSrv.GracefulStop()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

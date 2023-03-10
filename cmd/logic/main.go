package main

import (
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/internal/logic/apigrpc"
	"github.com/zhixunjie/im-fun/internal/logic/apihttp"
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
	httpSrv := apihttp.New(conf.Conf, svc)
	// init GRPC server
	rpcSrv := apigrpc.New(conf.Conf.RPC.Server, svc)

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

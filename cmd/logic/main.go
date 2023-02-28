package main

import (
	log "github.com/golang/glog"
	"github.com/zhixunjie/im-fun/internal/logic/apigrpc"
	"github.com/zhixunjie/im-fun/internal/logic/apihttp"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/internal/logic/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
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
		log.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			httpSrv.Close()
			rpcSrv.GracefulStop()
			log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

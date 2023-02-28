package main

import (
	"github.com/golang/glog"
	"github.com/zhixunjie/im-fun/internal/comet"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/internal/comet/grpc"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	if err := conf.InitConfig(); err != nil {
		panic(err)
	}
	InitCommon()
	// init server
	srv := comet.NewServer(conf.Conf)
	// init TCP server
	// init WS server
	// init GRPC server
	rpcSrv := grpc.New(srv, conf.Conf.RPC.Server)

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		glog.Infof("get signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			rpcSrv.GracefulStop()
			srv.Close()
			glog.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func InitCommon() {
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func InitDiscovery() {
	// TBD
}

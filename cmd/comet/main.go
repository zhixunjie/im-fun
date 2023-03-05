package main

import (
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/internal/comet"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/internal/comet/grpc"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	var err error
	if err = conf.InitConfig(); err != nil {
		panic(err)
	}
	InitCommon()
	// init server
	srv := comet.NewServer(conf.Conf)
	// init TCP server
	var lisTCPSrv *net.TCPListener
	if lisTCPSrv, err = comet.InitTCP(srv, runtime.NumCPU()); err != nil {
		panic(err)
	}
	// init WS server
	var lisWebSocketSrv *net.TCPListener
	if lisWebSocketSrv, err = comet.InitWs(srv, runtime.NumCPU()); err != nil {
		panic(err)
	}
	// init GRPC server
	rpcSrv := grpc.New(srv, conf.Conf.RPC.Server)

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logrus.Infof("get signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			rpcSrv.GracefulStop()
			lisTCPSrv.Close()
			lisWebSocketSrv.Close()
			srv.Close()
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

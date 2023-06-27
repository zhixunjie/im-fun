package main

import (
	"github.com/zhixunjie/im-fun/internal/comet"
	commetgrpc "github.com/zhixunjie/im-fun/internal/comet/api/grpc"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/perf"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	// init pprof
	perf.InitPProf("127.0.0.1:6060")
	// init config
	var err error
	if err = conf.InitConfig("cmd/comet/comet.yaml"); err != nil {
		panic(err)
	}
	// init common
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	// init server
	var lis1, lis2 *net.TCPListener
	srv := comet.NewServer(conf.Conf)
	{
		// init TCP server
		if lis1, err = comet.InitTCP(srv, runtime.NumCPU(), channel.ConnTypeTcp); err != nil {
			panic(err)
		}
		// init WebSocket server
		if lis2, err = comet.InitTCP(srv, runtime.NumCPU(), channel.ConnTypeWebSocket); err != nil {
			panic(err)
		}
	}

	// init GRPC server
	rpcSrv := commetgrpc.New(srv, conf.Conf.RPC.Server)

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logging.Infof("get signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			rpcSrv.GracefulStop()
			lis2.Close()
			lis1.Close()
			srv.Close()
			return
		case syscall.SIGHUP:
			return
		default:
			return
		}
	}
}

func InitDiscovery() {
	// TBD
}

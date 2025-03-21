package main

import (
	"context"
	"flag"
	"github.com/zhixunjie/im-fun/internal/comet"
	commetgrpc "github.com/zhixunjie/im-fun/internal/comet/api/grpc"
	"github.com/zhixunjie/im-fun/internal/comet/channel"
	"github.com/zhixunjie/im-fun/internal/comet/conf"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/perf"
	"github.com/zhixunjie/im-fun/pkg/prometheus/register"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var (
	confPath string
)

// get params
func init() {
	flag.StringVar(&confPath, "conf", "cmd/comet/comet.yaml", "配置文件的路径")
	flag.Parse()
}

func main() {
	// init pprof
	perf.InitPProf(":6060")
	// init prometheus
	register.InitProm(":7060")
	// init config
	var err error
	if err = conf.InitConfig(confPath); err != nil {
		panic(err)
	}
	// init common
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	// init server
	var lis1, lis2 *net.TCPListener
	srv, instance := comet.NewTcpServer(conf.Conf)
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

	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// init GRPC server
	rpcSrv, deRegister, err := commetgrpc.NewServer(ctx, srv, conf.Conf, instance)
	if err != nil {
		panic(err)
	}

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logging.Infof("get signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			deRegister()
			rpcSrv.GracefulStop()
			_ = lis2.Close()
			_ = lis1.Close()
			_ = srv.Close()
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

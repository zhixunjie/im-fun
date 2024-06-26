package main

import (
	"flag"
	"github.com/zhixunjie/im-fun/internal/job"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/perf"
	"github.com/zhixunjie/im-fun/pkg/prometheus/register"
	"os"
	"os/signal"
	"syscall"
)

var (
	confPath string
)

// get params
func init() {
	flag.StringVar(&confPath, "conf", "cmd/job/job.yaml", "配置文件的路径")
	flag.Parse()
}

func main() {
	// init pprof
	perf.InitPProf(":6061")
	// init prometheus
	register.InitProm(":7061")
	// init config
	if err := conf.InitConfig(confPath); err != nil {
		panic(err)
	}
	j := job.NewJob(conf.Conf)

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logging.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			j.Close()
			return
		case syscall.SIGHUP:
			return
		default:
			return
		}
	}
}

package main

import (
	"github.com/zhixunjie/im-fun/internal/job"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/perf"
	"github.com/zhixunjie/im-fun/pkg/prometheus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// init pprof
	perf.InitPProf(":6061")
	// init prometheus
	prometheus.InitPrometheus(":7061")
	// init config
	if err := conf.InitConfig("cmd/job/job.yaml"); err != nil {
		panic(err)
	}
	j := job.New(conf.Conf)

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

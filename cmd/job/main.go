package main

import (
	"github.com/zhixunjie/im-fun/internal/job"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/perf"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// init pprof
	perf.InitPProf("127.0.0.1:6061")
	// init config
	if err := conf.InitConfig("cmd/job/job.yaml"); err != nil {
		panic(err)
	}
	job := job.New(conf.Conf)

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logging.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			job.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

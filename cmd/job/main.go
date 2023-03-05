package main

import (
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/internal/job"
	"github.com/zhixunjie/im-fun/internal/job/conf"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := conf.InitConfig(); err != nil {
		panic(err)
	}
	job := job.New(conf.Conf)

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logrus.Infof("get a signal %s", s.String())
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

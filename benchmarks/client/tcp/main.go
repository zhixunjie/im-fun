package main

import (
	"flag"
	"fmt"
	"github.com/zhixunjie/im-fun/benchmarks/client/tcp/operation"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"math/rand"
	"runtime"
	"time"
)

var (
	addr string
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())

	// get params
	var start, num int64
	flag.Int64Var(&start, "start", 0, "用户ID的开始值")
	flag.Int64Var(&num, "num", 0, "启动X个客户端")
	flag.StringVar(&addr, "addr", "", "服务端地址")
	flag.Parse()

	// check params
	if addr == "" {
		fmt.Printf("没有指定参数 addr")
		return
	}
	if start == 0 || num == 0 {
		fmt.Printf("start或num参数等于0")
		return
	}

	// start to run
	go operation.DashBoard()
	var i int64
	for i = start; i < start+num; i++ {
		go func(userId int64) {
			for {
				// 切分QPS
				sec := rand.Intn(120)
				logging.Infof("userId=%v try to connect server after %v second", userId, sec)
				time.Sleep(time.Duration(sec) * time.Second)
				// start
				operation.Start(userId, addr)

				// restart after some second
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
			}
		}(i)
	}

	// signal
	var exit chan bool
	<-exit
}

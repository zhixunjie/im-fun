package main

import (
	"flag"
	"fmt"
	"github.com/zhixunjie/im-fun/benchmarks/client/tcp/operation"
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
	flag.Int64Var(&num, "num", 0, "客户端的数据")
	flag.StringVar(&addr, "addr", "", "服务端地址")
	flag.Parse()

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
				operation.Start(userId, addr)
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
			}
		}(i)
	}

	// signal
	var exit chan bool
	<-exit
}

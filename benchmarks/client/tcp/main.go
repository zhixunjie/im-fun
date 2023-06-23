package main

// Start Commond eg: ./client 1 1000 localhost:3101
// first parameter：beginning userId
// second parameter: amount of clients
// third parameter: comet server ip

import (
	"flag"
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

	// get params
	var begin, num int64
	flag.Int64Var(&begin, "begin", 0, "用户ID的开始值")
	flag.Int64Var(&num, "num", 0, "客户端的数据")
	flag.StringVar(&addr, "addr", "", "服务端地址")
	flag.Parse()

	// begin to run
	go operation.DashBoard()
	var i int64
	for i = begin; i < begin+num; i++ {
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

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zhixunjie/im-fun/benchmarks/send"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/pkg/http"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	addr   string
	client *http.Client
)

func init() {
	client = http.NewClient(http.RequestTimeout(time.Second * 3))
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())

	// get params
	var start, num int64
	var duration int
	flag.Int64Var(&start, "start", 0, "用户ID的开始值")
	flag.Int64Var(&num, "num", 0, "发送消息的次数")
	flag.StringVar(&addr, "addr", "", "服务端地址")
	flag.IntVar(&duration, "duration", 0, "持续时间")
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

	// calc params
	cpuNum := runtime.NumCPU() * 2
	step := num / int64(cpuNum)

	// set timer
	timer := time.NewTimer(time.Duration(duration) * time.Second)
	go func() {
		select {
		case <-timer.C:
			os.Exit(0)
		}
	}()

	// begin to run
	var wg sync.WaitGroup
	st, ed := start, start+step
	for i := 0; i < cpuNum; i++ {
		wg.Add(1)
		go func(s, d int64) {
			//fmt.Println(s, d)
			defer wg.Done()
			Start(s, d)
		}(st, ed)
		st += step
		ed += step
	}
	wg.Wait()
}

func Start(st, ed int64) {
	userId := send.RandUserId()
	logging.Infof("send start[%d~%d],userId=%v", st, ed, userId)

	// build msg
	msg := request.SendToUsersByIdsReq{
		UserIds: []uint64{userId},
		Message: send.Msg,
	}

	// build body
	bodyStr, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	// start to request
	var successCount, failCount int
	var rsp http.Response
	reqUrl := fmt.Sprintf("http://%v%v", addr, send.UrlSendToUsersByIds)
	for i := st; i < ed; i++ {
		rsp, err = client.Post(reqUrl, bytes.NewBuffer(bodyStr), nil, nil)
		if err != nil {
			failCount++
			continue
		}
		successCount++
		logging.Infof("rsp=%s", rsp.Body)
		time.Sleep(50 * time.Millisecond)
	}
	logging.Infof("send end[%d~%d],userId=%v,success=%d,fail=%d", st, ed, userId, successCount, failCount)
}

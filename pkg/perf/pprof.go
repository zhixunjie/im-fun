package perf

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func InitPProf(addr string) {
	// 启动pprof的HTTP服务器
	go func() {
		fmt.Println("start pprof HTTP Server")
		_ = http.ListenAndServe(addr, nil)
	}()
}

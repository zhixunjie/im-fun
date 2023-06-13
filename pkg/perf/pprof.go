package perf

import (
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net/http"
	_ "net/http/pprof"
)

func InitPProf(addr string) {
	// 启动pprof的HTTP服务器
	go func() {
		logging.Infof("start pprof HTTP Server")
		_ = http.ListenAndServe(addr, nil)
	}()
}

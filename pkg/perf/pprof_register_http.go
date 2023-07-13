package perf

import (
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net/http"
	_ "net/http/pprof"
)

// InitPProf 使用新的多路复用器进行监听
func InitPProf(addr string) {
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// 启动pprof的HTTP服务器
	go func() {
		logging.Infof("start pprof HTTP Server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Errorf("listen: %s,err=%v", addr, err)
			return
		}
	}()
}

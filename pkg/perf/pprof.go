package register

import (
	"github.com/felixge/fgprof"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net/http"
	"net/http/pprof"
)

// InitPProf 使用新的多路复用器进行监听
func InitPProf(addr string) {
	// new mux
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	// use fgprof:
	// note: Please upgrade to Go 1.19 or newer. In older versions of Go fgprof can cause significant STW latencies in applications with a lot of goroutines (> 1-10k).
	// See CL 387415 for more details.
	mux.Handle("/debug/fgprof", fgprof.Handler())

	// bind mux to server
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

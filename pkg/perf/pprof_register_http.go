package register

import (
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net/http"
	"net/http/pprof"
)

/**
为什么引入net/http/pprof后，就能使用默认的多路复用器进行监听？
因为该包的init函数执行了如下的代码：

func init() {
	http.HandleFunc("/debug/pprof/", Index)
	http.HandleFunc("/debug/pprof/cmdline", Cmdline)
	http.HandleFunc("/debug/pprof/profile", Profile)
	http.HandleFunc("/debug/pprof/symbol", Symbol)
	http.HandleFunc("/debug/pprof/trace", Trace)
}
*/

func InitPProf1(addr string) {
	// 启动pprof的HTTP服务器
	go func() {
		logging.Infof("start pprof HTTP Server")
		if err := http.ListenAndServe(addr, nil); err != nil && err != http.ErrServerClosed {
			logging.Errorf("listen: %s,err=%v", addr, err)
			return
		}
	}()
}

// InitPProf 使用新的多路复用器进行监听
func InitPProf(addr string) {
	// new mux
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

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

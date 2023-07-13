package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net/http"
)

// InitPrometheus 实现Go应用
// - https://hulining.gitbook.io/prometheus/guides/go-application
// - https://prometheus.io/docs/guides/go-application/
func InitPrometheus(addr string) {
	// new mux
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	// bind mux to server
	srv := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// 启动prometheus的HTTP服务器
	go func() {
		logging.Infof("start Prometheus HTTP Server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Errorf("listen: %s,err=%v", addr, err)
			return
		}
	}()
}

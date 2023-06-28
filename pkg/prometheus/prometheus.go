package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net/http"
)

func InitPrometheus(addr string) {
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// 启动prometheus的HTTP服务器
	mux.Handle("/metrics", promhttp.Handler())
	go func() {
		logging.Infof("start Prometheus HTTP Server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Errorf("listen: %s,err=%v", addr, err)
			return
		}
	}()
}

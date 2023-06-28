package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net/http"
)

func InitPrometheus(addr string) {
	// 启动prometheus的HTTP服务器
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		logging.Infof("start pprof HTTP Server")
		_ = http.ListenAndServe(addr, nil)
	}()
}

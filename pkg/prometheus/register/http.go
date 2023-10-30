package register

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net/http"
)

// 使用新定义的注册中心（正式项目中使用）

func InitProm(addr string) {
	// new mux
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(M.G, promhttp.HandlerOpts{Registry: M.R}))

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

package http

import (
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net/http"
)

// 使用net/http的默认多路器

// DefaultMux 默认的多路复用器
// 测试： curl http://127.0.0.1:8080/?a=1&b=2
func DefaultMux() {
	// 设置默认多路复用器的路由信息
	http.HandleFunc("/ws", WsHandler)

	// nil表示使用默认的多路复用器(DefaultServeMux)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logging.Errorf("创建监听服务器失败: err=%v", err)
		return
	}
}

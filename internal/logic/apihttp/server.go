package apihttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/internal/logic/dao"
	"github.com/zhixunjie/im-fun/internal/logic/service"
)

// Server is http server.
type Server struct {
	engine *gin.Engine
	svc    *service.Service
}

func New(conf *conf.Config, svc *service.Service) *Server {
	dao.InitDao()
	engine := gin.Default()

	srv := &Server{
		engine: engine,
		svc:    svc,
	}
	// 设置-路由
	srv.SetupRouter()

	// begin to listen
	go func() {
		fmt.Printf("HTTP服务启动成功，正在监听：%v\n", conf.HTTPServer.Addr)
		if err := engine.Run(conf.HTTPServer.Addr); err != nil {
			panic(err)
		}
	}()
	return srv
}

func (s *Server) Close() {

}

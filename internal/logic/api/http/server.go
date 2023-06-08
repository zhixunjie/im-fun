package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/internal/logic/service"
	"net/http"
	"time"
)

// Server is http server.
type Server struct {
	engine     *gin.Engine
	httpServer *http.Server
	svc        *service.Service
}

func New(conf *conf.Config, svc *service.Service) *Server {
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()

	srv := &Server{
		engine: engine,
		svc:    svc,
	}
	// 设置-路由
	srv.SetupRouter()

	// set net.http
	addr := conf.HTTPServer.Addr
	httpServer := &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	srv.httpServer = httpServer

	// begin to listen
	logrus.Infof("HTTP server is listening：%v", addr)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("ListenAndServe,err=%v", err)
		}
	}()
	return srv
}

func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// invoke Shutdown：stop accept new connection && will be forced to stop after 5 seconds
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server Shutdown Failed:%+v", err)
	}
	logrus.Infof("Server Exited Properly")
}

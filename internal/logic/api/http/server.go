package http

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"github.com/zhixunjie/im-fun/internal/logic/biz"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net/http"
	"time"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServer)

// Server is http server.
type Server struct {
	engine     *gin.Engine
	httpServer *http.Server

	bz             *biz.Biz
	BzContact      *biz.ContactUseCase
	BzMessage      *biz.MessageUseCase
	BzGroupMessage *biz.GroupMessageUseCase
	BzUser         *biz.UserUseCase
	BzUserGroup    *biz.UserGroupUseCase
}

func NewServer(conf *conf.Config, bz *biz.Biz,
	bzContact *biz.ContactUseCase, bzMessage *biz.MessageUseCase, bzGroupMessage *biz.GroupMessageUseCase,
	bzUser *biz.UserUseCase, bzUserGroup *biz.UserGroupUseCase,
) *Server {
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// get gin engine
	engine := gin.Default()

	// set net.http
	addr := conf.HTTPServer.Addr
	httpServer := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	// set Server
	srv := &Server{
		engine:         engine,
		httpServer:     httpServer,
		bz:             bz,
		BzContact:      bzContact,
		BzMessage:      bzMessage,
		BzGroupMessage: bzGroupMessage,
		BzUser:         bzUser,
		BzUserGroup:    bzUserGroup,
	}
	// 设置-中间件
	srv.SetupMiddleware()
	// 设置-路由
	srv.SetupRouter()

	// begin to listen
	logging.Infof("HTTP server is listening：%v", addr)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	logging.Infof("Server Exited Properly")
}

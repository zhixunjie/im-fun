package response

import (
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"net/http"
)

type Base struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func JsonError(ctx *gin.Context, err error) {
	logging.Errorf("method=%v,err=%v", ctx.Request.Method, err)
	ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "msg": "error: " + err.Error()})
	return
}

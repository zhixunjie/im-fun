package response

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Base struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func JsonError(ctx *gin.Context, err error) {
	logrus.Errorf("method=%v,err=%v", ctx.Request.Method, err)
	ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "msg": "error: " + err.Error()})
	return
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhixunjie/im-fun/internal/logic/api/http"
	"github.com/zhixunjie/im-fun/internal/logic/dao"
)

func main() {
	dao.InitDao()
	engine := gin.Default()

	// 设置-路由
	http.SetupRouter(engine)

	// 开始执行
	if err := engine.Run(":8080"); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}
}

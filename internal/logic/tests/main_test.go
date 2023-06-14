package tests

import (
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/internal/logic/service"
	"os"
	"testing"
)

var GlobalSvc *service.Service

func TestMain(main *testing.M) {
	// 加载本地文件配置
	if err := conf.InitConfig("cmd/logic/logic.yaml"); err != nil {
		panic(err)
	}
	GlobalSvc = service.New(conf.Conf)

	// run testing(执行测试用例)
	main.Run()

	// after testing
	fmt.Println("testing finish")

	os.Exit(0)
}

package tests

import (
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"os"
	"testing"
)

var GlobalSvc *biz.Service

func TestMain(main *testing.M) {
	// 加载本地文件配置
	if err := conf.InitConfig("../../../cmd/logic/logic.yaml"); err != nil {
		panic(err)
	}
	GlobalSvc = biz.New(conf.Conf)

	// run testing(执行测试用例)
	main.Run()

	// after testing
	fmt.Println("testing finish")

	os.Exit(0)
}

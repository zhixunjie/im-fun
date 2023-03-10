package tests

import (
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/dao"
	"os"
	"testing"
)

func TestMain(main *testing.M) {
	// 加载本地文件配置
	dao.InitDao()

	// run testing(执行测试用例)
	main.Run()

	// after testing
	fmt.Println("testing finish")

	os.Exit(0)
}

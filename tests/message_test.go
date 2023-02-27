package tests

import (
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/dao"
	"testing"
)

func TestMessage1(t *testing.T) {
	res, _ := dao.QueryMsgLogic(1001)
	fmt.Printf("%+v\n", res)
}

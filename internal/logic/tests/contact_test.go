package tests

import (
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/dao"
	"testing"
)

func TestContact1(t *testing.T) {
	res, _ := dao.QueryContactLogic(1001, 1002)
	fmt.Printf("%+v\n", res)
}

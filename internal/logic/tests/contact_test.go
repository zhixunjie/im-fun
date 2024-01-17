package tests

import (
	"fmt"
	"github.com/zhixunjie/im-fun/cmd/logic/wire"
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"log"
	"testing"
)

func TestQueryContactLogic(t *testing.T) {
	res, err := wire.GetContactRepo(conf.Conf).QueryContactLogic(1001, 1002)
	if err != nil {
		log.Fatal()
	}
	fmt.Printf("%+v\n", res)
}

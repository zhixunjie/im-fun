package tests

import (
	"fmt"
	"testing"
)

func TestMessage1(t *testing.T) {
	res, _ := GlobalSvc.GetDao().QueryMsgLogic(1001)
	fmt.Printf("%+v\n", res)
}

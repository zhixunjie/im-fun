package tests

import (
	"fmt"
	"testing"
)

func TestContact1(t *testing.T) {
	res, _ := GlobalSvc.GetDao().QueryContactLogic(1001, 1002)
	fmt.Printf("%+v\n", res)
}

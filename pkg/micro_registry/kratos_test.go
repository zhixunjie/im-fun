package micro_registry

import (
	"fmt"
	"net"
	"testing"
)

func TestBuildInstance(t *testing.T) {
	addr := ":8123"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(BuildInstance("comet", addr, lis))
}

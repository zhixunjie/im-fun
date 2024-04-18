package registry

import (
	"fmt"
	"testing"
)

func TestBuildInstance(t *testing.T) {
	addr := ":8123"
	fmt.Println(BuildServiceInstance("comet", "tcp", addr))
	fmt.Println(BuildServiceInstance("logic", "tcp", addr))
}

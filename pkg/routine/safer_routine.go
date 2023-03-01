package routine

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

func Go(ctx context.Context, f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				err := fmt.Errorf("goroutine: panic recovered: %s\n%s", r, buf)
				logrus.Errorf("goroutine panic: %s", err)
			}
		}()
		f()
	}()
}

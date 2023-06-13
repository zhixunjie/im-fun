package routine

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"runtime"
)

func Go(ctx context.Context, f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				err := fmt.Errorf("goroutine: panic recovered: %s\n%s", r, buf)
				logging.Errorf("goroutine panic: %s", err)
			}
		}()
		f()
	}()
}

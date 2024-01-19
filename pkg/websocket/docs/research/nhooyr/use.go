package nhooyr

import (
	"context"
	"fmt"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

// https://github.com/nhooyr/websocket
// 记录关键节点：整个跳转索引就行（比较容易看懂）
func nhooyr() {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			panic(err)
		}
		defer c.Close(websocket.StatusInternalError, "the sky is falling")

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()

		c.Reader(ctx)
		// c.WriterPool(ctx)
		// var v interface{}
		// wsjson.Read(ctx, c, &v)
		// wsjson.Write(ctx, c, &v)

		c.Close(websocket.StatusNormalClosure, "")
	})
	fmt.Println(handlerFunc)
}

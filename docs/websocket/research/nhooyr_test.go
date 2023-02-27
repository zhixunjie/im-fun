package research

import (
	"context"
	"fmt"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"testing"
	"time"
)

// https://github.com/nhooyr/websocket
// 整个跳转索引就行
func TestServer11(t *testing.T) {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			panic(err)
		}
		defer c.Close(websocket.StatusInternalError, "the sky is falling")

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()

		var v interface{}
		wsjson.Read(ctx, c, &v)
		wsjson.Write(ctx, c, &v)

		c.Close(websocket.StatusNormalClosure, "")
	})
	fmt.Println(handlerFunc)
}

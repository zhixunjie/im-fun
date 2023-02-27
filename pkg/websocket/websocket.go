package websocket

import (
	"context"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func handshake() {

	http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
		}
		defer c.Close(websocket.StatusInternalError, "the sky is falling")

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()

		var v interface{}
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			// ...
		}

		log.Printf("received: %v", v)

		c.Close(websocket.StatusNormalClosure, "")
	})
}

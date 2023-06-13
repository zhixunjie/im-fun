package research

import (
	"github.com/gorilla/websocket"
	"testing"
)

// https://github.com/gorilla/websocket
// 整个跳转索引就行
func TestServer1(t *testing.T) {
	// var upgrader websocket.Upgrader
	// upgrader.Upgrade()

	var conn *websocket.Conn
	conn.ReadMessage()
	// conn.WriteMessage()
	conn.Close()

}

package gorilla

import (
	"github.com/gorilla/websocket"
)

// https://github.com/gorilla/websocket
// 记录关键节点：整个跳转索引就行
func gorilla() {
	// var upgrader websocket.Upgrader
	// upgrader.Upgrade()

	var conn *websocket.Conn
	conn.ReadMessage()
	// conn.WriteMessage()
	conn.Close()
}

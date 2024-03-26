package channel

type ConnType int

const (
	ConnTypeTcp ConnType = iota + 1
	ConnTypeWebSocket
)

func LogHeadByConnType(connType ConnType) string {
	if connType == ConnTypeWebSocket {
		return "WebSocket|"
	}
	return "TCP|"
}

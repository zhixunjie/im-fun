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

type CleanPath int

const (
	CleanPath1 CleanPath = iota + 1
	CleanPath2
	CleanPath3
)

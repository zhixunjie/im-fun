package model

import "github.com/zhixunjie/im-fun/pkg/tcp"

// Proto proto.
type Proto struct {
	PackLen   int32  // package length
	HeaderLen int16  // header length
	Ver       int16  // protocol version
	Op        int32  // operation for request
	Seq       int32  // sequence number chosen by client
	Body      []byte // body
	BodyLen   int32  // body length
}

type AuthParams struct {
	TcpSessionId *tcp.SessionId `json:"tcp_session_id"` // 唯一地标识一条TCP连接
	RoomId       string         `json:"room_id"`
	Platform     string         `json:"platform"`
	Token        string         `json:"token"`
}

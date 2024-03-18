package channel

import (
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/pkg/tcp"
)

type AuthParams struct {
	UserInfo UserInfo `json:"user_info"`
	Token    string   `json:"token"`
}

type UserInfo struct {
	TcpSessionId *tcp.SessionId `json:"tcp_session_id"` // 唯一地标识一条TCP连接
	RoomId       string         `json:"room_id"`        // 房间ID
	Platform     pb.Platform    `json:"platform"`       // 客户端平台
	IP           string         `json:"ip"`             // 客户端的IP地址
}

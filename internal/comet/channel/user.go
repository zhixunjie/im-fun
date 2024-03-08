package channel

import (
	"github.com/zhixunjie/im-fun/api/pb"
	"github.com/zhixunjie/im-fun/pkg/tcp"
)

type AuthParams struct {
	UserInfo
	Token string `json:"token"`
}

type UserInfo struct {
	UserId       uint64         `json:"user_id"`        // 用户ID
	UserKey      string         `json:"user_key"`       // 用户KEY
	TcpSessionId *tcp.SessionId `json:"tcp_session_id"` // 唯一地标识一条TCP连接
	RoomId       string         `json:"room_id"`        // 房间ID
	Platform     pb.Platform    `json:"platform"`       // 客户端平台
	IP           string         `json:"ip"`             // 客户端的IP地址
}

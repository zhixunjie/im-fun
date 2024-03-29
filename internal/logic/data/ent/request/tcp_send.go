package request

import "github.com/zhixunjie/im-fun/pkg/tcp"

type SendToUsersReq struct {
	TcpSessionIds []tcp.SessionId `json:"tcp_session_ids"`
	SubId         int32           `json:"sub_id"`
	Message       string          `json:"message"`
}

type SendToUsersByIdsReq struct {
	UserIds []uint64 `json:"user_ids"`
	SubId   int32    `json:"sub_id"`
	Message string   `json:"message"`
}

type SendToRoomReq struct {
	RoomId   string `json:"room_id"`
	RoomType string `json:"room_type"`
	SubId    int32  `json:"sub_id"`
	Message  string `json:"message"`
}

type SendToAllReq struct {
	Speed   int32  `json:"speed"`
	SubId   int32  `json:"sub_id"`
	Message string `json:"message"`
}

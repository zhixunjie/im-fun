package request

import "github.com/zhixunjie/im-fun/api/pb"

type SendToUsersReq struct {
	Atom          *pb.Atom
	TcpSessionIds []string `json:"tcp_session_ids"`
	SubId         int32    `json:"sub_id"`
	Message       string   `json:"message"`
}

type SendToUsersByIdsReq struct {
	Atom    *pb.Atom
	UniIds  []string `json:"uni_ids"`
	SubId   int32    `json:"sub_id"`
	Message string   `json:"message"`
}

type SendToRoomReq struct {
	Atom     *pb.Atom
	RoomId   string `json:"room_id"`
	RoomType string `json:"room_type"`
	SubId    int32  `json:"sub_id"`
	Message  string `json:"message"`
}

type SendToAllReq struct {
	Atom    *pb.Atom
	Speed   int32  `json:"speed"`
	SubId   int32  `json:"sub_id"`
	Message string `json:"message"`
}

type OnlineUniIdReq struct {
	Atom *pb.Atom
}

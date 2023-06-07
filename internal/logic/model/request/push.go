package request

type SendToUserKeysReq struct {
	UserKeys []string `json:"user_keys"`
	SubId    int32    `json:"sub_id"`
	Message  string   `json:"message"`
}

type SendToUserIdsReq struct {
	UserIds []int64 `json:"user_ids"`
	SubId   int32   `json:"sub_id"`
	Message []byte  `json:"message"`
}

type SendToRoomReq struct {
	RoomId   string `json:"room_id"`
	RoomType string `json:"room_type"`
	SubId    int32  `json:"sub_id"`
	Message  []byte `json:"message"`
}

type SendToAllReq struct {
	Speed   int32  `json:"speed"`
	SubId   int32  `json:"sub_id"`
	Message []byte `json:"message"`
}

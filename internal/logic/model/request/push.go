package request

type PushUserKeysReq struct {
	UserKeys []string `json:"user_keys"`
	SubId    int32    `json:"sub_id"`
	Message  []byte   `json:"message"`
}

type PushUserIdsReq struct {
	UserIds []int64 `json:"user_ids"`
	SubId   int32   `json:"sub_id"`
	Message []byte  `json:"message"`
}

type PushUserRoomReq struct {
	RoomId   string `json:"room_id"`
	RoomType string `json:"room_type"`
	SubId    int32  `json:"sub_id"`
	Message  []byte `json:"message"`
}

type PushUserAllReq struct {
	Speed   int32  `json:"speed"`
	SubId   int32  `json:"sub_id"`
	Message []byte `json:"message"`
}

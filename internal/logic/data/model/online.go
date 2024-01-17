package model

type Online struct {
	Server    string           `json:"server"`
	RoomCount map[string]int32 `json:"room_count"`
	Updated   int64            `json:"updated"`
}

type Top struct {
	RoomID string `json:"room_id"`
	Count  int32  `json:"count"`
}

package model

// Proto proto.
type Proto struct {
	PackLen   int32  // package length
	HeaderLen int16  // header length
	BodyLen   int32  // body length
	Ver       int16  // protocol version
	Op        int32  // operation for request
	Seq       int32  // sequence number chosen by client
	Body      []byte // body
}

type AuthParams struct {
	UserId   int64  `json:"user_id"`
	UserKey  string `json:"user_key"`
	RoomId   string `json:"room_id"`
	Platform string `json:"platform"`
	Token    string `json:"token"`
}

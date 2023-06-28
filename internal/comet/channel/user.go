package channel

type AuthParams struct {
	UserInfo
	Token string `json:"token"`
}

type UserInfo struct {
	UserId   int64  `json:"user_id"`  // 用户ID
	UserKey  string `json:"user_key"` // 用户KEY
	RoomId   string `json:"room_id"`  // 房间ID
	Platform string `json:"platform"` // 平台（客户端（PC、安卓、IOS）、Web、小程序等）
	IP       string `json:"ip"`
}

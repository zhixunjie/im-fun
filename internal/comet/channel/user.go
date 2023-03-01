package channel

type UserInfo struct {
	UserId   int64  // 用户ID
	UserKey  string // 用户KEY
	RoomId   string // 房间ID
	Platform string // 平台（客户端（PC、安卓、IOS）、Web、小程序等）
	IP       string
}

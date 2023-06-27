package send

import "math/rand"

const Msg = "Hello World"

const (
	UrlSendUserKeys = "/send/user/keys"
	UrlSendUserIds  = "/send/user/ids"
	UrlSendRoom     = "/send/user/room"
	UrlSendAll      = "/send/user/all"
)

const MaxUserId = 1000

// RandUserId 随机一个UserId
// value scope：[0,MaxUserId)
func RandUserId() int64 {
	return rand.Int63n(MaxUserId)
}

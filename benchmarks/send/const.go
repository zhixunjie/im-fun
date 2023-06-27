package send

import "math/rand"

const Msg = "Hello World"

const (
	UrlSendUserKeys = "/im/send/user/keys"
	UrlSendUserIds  = "/im/send/user/ids"
	UrlSendRoom     = "/im/send/user/room"
	UrlSendAll      = "/im/send/user/all"
)

const MaxUserId = 1000

// RandUserId 随机一个UserId
// value scope：[0,MaxUserId)
func RandUserId() int64 {
	return rand.Int63n(MaxUserId)
}

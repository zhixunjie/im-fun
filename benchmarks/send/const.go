package send

import (
	"github.com/spf13/cast"
	"math/rand"
)

const Msg = "Hello World"

const (
	UrlSendToUsers      = "/im/send/to/users"
	UrlSendToUsersByIds = "/im/send/to/users/by/ids"
	UrlSendToRoom       = "/im/send/to/room"
	UrlSendToAll        = "/im/send/to/all"
)

const MaxUserId = 1000

// RandUserId 随机一个UserId
// value scope：[0,MaxUserId)
func RandUserId() uint64 {
	return cast.ToUint64(rand.Int63n(MaxUserId))
}

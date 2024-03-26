package tcp

import "fmt"

// SessionId 唯一地标识一条TCP连接
type SessionId struct {
	UserId  uint64 `json:"user_id"`  // 唯一ID，一般使用用户ID
	UserKey string `json:"user_key"` // ID维度下，作进一步的划分（比如：某个终端下的用户ID）
}

func (t *SessionId) ToString() string {
	return fmt.Sprintf("%v:%v", t.UserId, t.UserKey)
}

//func NewSessionId(userId uint64, userKey string) *SessionId {
//	return &SessionId{
//		UserId:  userId,
//		UserKey: userKey,
//	}
//}

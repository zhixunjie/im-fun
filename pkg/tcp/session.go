package tcp

import "fmt"

// SessionId 唯一地标识一条TCP连接
type SessionId struct {
	id     uint64 // 唯一ID，一般使用用户ID
	idType string // ID维度下，作进一步的划分（比如：某个终端下的用户ID）
}

func (t *SessionId) ToString() string {
	return fmt.Sprintf("%v:%v", t.id, t.idType)
}

func (t *SessionId) Id() uint64 {
	return t.id
}

func (t *SessionId) Type() string {
	return t.idType
}

func NewSessionId(id uint64, idType string) *SessionId {
	return &SessionId{
		id:     id,
		idType: idType,
	}
}

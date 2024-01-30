package request

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
)

// ContactFetchReq 拉取会话列表（by version_id）
type ContactFetchReq struct {
	VersionId model.BigIntType     `json:"version_id"` // 版本id
	OwnerId   model.BigIntType     `json:"owner_id"`   // 会话拥有者
	OwnerType gen_id.ContactIdType `json:"owner_type"` // 会话拥有者的用户类型
}

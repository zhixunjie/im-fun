package request

import (
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
)

// ContactFetchReq 拉取会话列表（by version_id）
type ContactFetchReq struct {
	VersionId model.BigIntType    `json:"version_id,string"` // 版本id
	Owner     *gen_id.ComponentId `json:"owner"`             // 会话拥有者
}

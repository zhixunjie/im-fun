package gen_id

import (
	"fmt"
	"github.com/zhixunjie/im-fun/pkg/utils"
)

// GetSessionId 获取会话id，小的uid在前，大的uid在后
func GetSessionId(uid1 uint64, uid2 uint64) string {
	smallerId, largerId := utils.GetSortNum(uid1, uid2)

	// session_id的组成部分：[ smallerId ":" largerId]
	return fmt.Sprintf("%d:%d", smallerId, largerId)
}

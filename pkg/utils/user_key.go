package utils

import "fmt"

// UserKeyComponent 标识一条TCP连接的ID（用户id + 用户终端key）
func UserKeyComponent(userId int64, userKey string) string {
	return fmt.Sprintf("%v_%v", userId, userKey)
}

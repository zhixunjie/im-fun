package utils

import "fmt"

func GetMergeUserKey(userId int64, userKey string) string {
	return fmt.Sprintf("%v_%v", userId, userKey)
}

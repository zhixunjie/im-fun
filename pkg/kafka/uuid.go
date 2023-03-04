package kafka

import (
	uuid "github.com/satori/go.uuid"
	"strings"
)

func Uuid() string {
	u := uuid.NewV4()
	uuidStr := strings.Replace(u.String(), "-", "", -1)
	return uuidStr
}

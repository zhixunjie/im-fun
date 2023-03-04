package utils

import (
	"fmt"
	"net/url"
)

// EncodeRoomKey encode a room key.
func EncodeRoomKey(roomType string, roomId string) string {
	return fmt.Sprintf("%s://%s", roomType, roomId)
}

// DecodeRoomKey decode room key.
func DecodeRoomKey(key string) (string, string, error) {
	u, err := url.Parse(key)
	if err != nil {
		return "", "", err
	}
	return u.Scheme, u.Host, nil
}

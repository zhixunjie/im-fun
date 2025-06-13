package gmodel

import "github.com/golang-jwt/jwt/v5"

type Atom struct {
	SeqId    string   `query:"seq_id"`   // 序列号
	UID      int64    `query:"uid"`      // 用户id
	DID      string   `query:"did"`      // 设备id
	APP      APP      `query:"app"`      // app id
	Platform Platform `query:"platform"` // 平台
}

type (
	AuthClaims struct {
		jwt.RegisteredClaims
		*AuthUserInfo
	}
	AuthUserInfo struct {
		Uid uint64 `json:"uid"`
	}
)

type Platform string

const (
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
)

type APP string

const (
	AppIM = "im"
)

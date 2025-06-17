package gmodel

import "github.com/golang-jwt/jwt/v5"

type (
	AuthClaims struct {
		jwt.RegisteredClaims
		*AuthUserInfo
	}
	AuthUserInfo struct {
		Uid uint64 `json:"uid"`
	}
)

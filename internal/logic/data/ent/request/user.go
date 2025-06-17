package request

import "github.com/zhixunjie/im-fun/pkg/gmodel"

type LoginReq struct {
	gmodel.Atom
	AccountType    gmodel.AccountType `json:"account_type" binding:"required,gte=1,lte=2"` // 账号类型
	AccountID      string             `json:"account_id" binding:"required,gte=1"`         // 设备号/手机号码/三方平台账号
	ThirdLoginInfo *ThirdLoginInfo    `json:"third_login_info,omitempty"`                  // 第三方平台登录信息
}

type ThirdLoginInfo struct {
	AccessToken string `json:"access_token"` // access_token
	ExpireIn    int64  `json:"expire_in"`    // 过期时间戳
}

type RefreshTokenReq struct {
	gmodel.Atom
}

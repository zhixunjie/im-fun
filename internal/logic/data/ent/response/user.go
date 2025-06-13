package response

import "github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"

type (
	LoginRsp struct {
		Base
		Data *LoginData `json:"data"`
	}
	LoginData struct {
		User  *model.User `json:"user"`
		Token string      `json:"token"`
	}
)

type (
	RefreshTokenRsp struct {
		Base
		Data *RefreshTokenData `json:"data"`
	}
	RefreshTokenData struct {
		Token string `json:"token"`
	}
)

package response

import "github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"

type (
	SendToUsersResp struct {
		Base
	}
)

type (
	SendToUsersByIdsResp struct {
		Base
	}
)

type (
	SendToRoomResp struct {
		Base
	}
)

type (
	SendToAllResp struct {
		Base
	}
)
type (
	OnlineUniIdRsp struct {
		Base
		Data *OnlineUniIdRspData `json:"data"`
	}
	OnlineUniIdRspData struct {
		Users []*model.User `json:"users"`
	}
)

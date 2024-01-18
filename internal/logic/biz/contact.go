package biz

import (
	"context"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"math"
	"sort"
)

type ContactUseCase struct {
	contactRepo *data.ContactRepo
}

func NewContactUseCase(contactRepo *data.ContactRepo) *ContactUseCase {
	return &ContactUseCase{contactRepo: contactRepo}
}

func (b *ContactUseCase) FetchSession(ctx context.Context, req *request.FetchSessionReq) (resp response.FetchSessionResp, err error) {
	//logHead := "FetchSession|"
	limit := 50

	// 会话只会拉取最新的
	list, err := b.contactRepo.RangeList(&model.QueryContactParams{
		FetchType:      model.FetchTypeForward,
		OwnerId:        req.OwnerId,
		PivotVersionId: req.VersionId,
		Limit:          limit,
	})
	if err != nil {
		return
	}

	// rebuild list
	var retList []*response.Contact
	minVersionId := uint64(math.MaxUint64)
	maxVersionId := uint64(0)
	for _, item := range list {
		minVersionId = utils.Min(minVersionId, item.VersionID)
		maxVersionId = utils.Max(maxVersionId, item.VersionID)

		// build message list
		retList = append(retList, &response.Contact{
			OwnerID:      item.OwnerID,
			PeerID:       item.PeerID,
			PeerType:     item.PeerType,
			PeerAck:      item.PeerAck,
			LastMsg:      nil,
			VersionID:    item.VersionID,
			SortKey:      item.SortKey,
			Status:       item.Status,
			Labels:       item.Labels,
			UnreadMsgNum: 0,
		})
	}

	// sort: 按照sort_key排序（从小到大）
	sort.Slice(retList, func(i, j int) bool {
		return retList[i].SortKey < retList[j].SortKey
	})

	resp.Data = response.FetchSessionData{
		ContactList:   retList,
		NextVersionId: maxVersionId,
		HasMore:       len(list) == limit,
	}

	return
}

package biz

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
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

// Fetch 拉取会话
func (b *ContactUseCase) Fetch(ctx context.Context, req *request.ContactFetchReq) (resp response.ContactFetchRsp, err error) {
	logHead := fmt.Sprintf("Fetch,req=%v", req)
	limit := 50

	// 会话只会拉取最新的
	ownerId := gen_id.NewComponentId(req.OwnerId, uint32(req.OwnerType))
	list, err := b.contactRepo.RangeList(logHead, &model.FetchContactRangeParams{
		FetchType:      model.FetchTypeForward,
		OwnerId:        ownerId,
		PivotVersionId: req.VersionId,
		Limit:          limit,
	})
	if err != nil {
		return
	}

	// rebuild list
	var retList []*response.ContactEntity
	minVersionId := uint64(math.MaxUint64)
	maxVersionId := uint64(0)
	for _, item := range list {
		minVersionId = utils.Min(minVersionId, item.VersionID)
		maxVersionId = utils.Max(maxVersionId, item.VersionID)

		// build message list
		retList = append(retList, &response.ContactEntity{
			OwnerID:      item.OwnerID,
			OwnerType:    gen_id.ContactIdType(item.OwnerType),
			PeerID:       item.PeerID,
			PeerType:     gen_id.ContactIdType(item.PeerType),
			PeerAck:      model.PeerAckStatus(item.PeerAck),
			VersionID:    item.VersionID,
			SortKey:      item.SortKey,
			Status:       model.ContactStatus(item.Status),
			Labels:       item.Labels,
			LastMsg:      nil,
			UnreadMsgNum: 0,
		})
	}

	// sort: 返回之前进行重新排序
	sort.Sort(response.ContactSortBySortKey(retList))

	resp.Data = response.FetchContactData{
		ContactList:   retList,
		NextVersionId: maxVersionId,
		HasMore:       len(list) == limit,
	}

	return
}

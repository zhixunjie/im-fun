package biz

import (
	"context"
	"fmt"
	"github.com/samber/lo"
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
	repoMessage *data.MessageRepo
}

func NewContactUseCase(contactRepo *data.ContactRepo, repoMessage *data.MessageRepo) *ContactUseCase {
	return &ContactUseCase{contactRepo: contactRepo, repoMessage: repoMessage}
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

	// extract: all peer ids
	peerIds := lo.Map(list, func(item *model.Contact, index int) *gen_id.ComponentId {
		return gen_id.NewComponentId(item.PeerID, item.PeerType)
	})
	retMap, err := b.repoMessage.MGetSessionUnread(ctx, logHead, ownerId, peerIds)
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
		peerId := gen_id.NewComponentId(item.PeerID, item.PeerType)

		var sessionUnreadCount int64
		if v, ok := retMap[peerId.ToString()]; ok {
			sessionUnreadCount = v
		}

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
			UnreadMsgNum: sessionUnreadCount,
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

package biz

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"sort"
)

type ContactUseCase struct {
	contactRepo *data.ContactRepo
	repoMessage *data.MessageRepo
}

func NewContactUseCase(contactRepo *data.ContactRepo, repoMessage *data.MessageRepo) *ContactUseCase {
	return &ContactUseCase{contactRepo: contactRepo, repoMessage: repoMessage}
}

// Fetch 拉取会话列表
func (b *ContactUseCase) Fetch(ctx context.Context, req *request.ContactFetchReq) (resp response.ContactFetchRsp, err error) {
	logHead := fmt.Sprintf("Fetch|req=%v", req)
	limit := 50

	// 会话只会拉取最新的
	ownerId := req.Owner
	list, err := b.contactRepo.RangeList(&model.FetchContactRangeParams{
		FetchType:      gmodel.FetchTypeForward,
		Owner:          ownerId,
		PivotVersionId: req.VersionId,
		Limit:          limit,
	})
	if err != nil {
		logging.Errorf(logHead+"RangeList err=%v", err)
		return
	}
	if len(list) == 0 {
		return
	}

	// extract: all peer ids
	peerIds := lo.Map(list, func(item *model.Contact, index int) *gmodel.ComponentId {
		return gmodel.NewComponentId(item.PeerID, gmodel.ContactIdType(item.PeerType))
	})
	retMap, err := b.repoMessage.MGetSessionUnread(ctx, logHead, ownerId, peerIds)
	if err != nil {
		return
	}

	// rebuild list
	var retList []*response.ContactEntity
	maxVersionId := uint64(0)
	for _, item := range list {
		maxVersionId = utils.Max(maxVersionId, item.VersionID)
		peerId := gmodel.NewComponentId(item.PeerID, gmodel.ContactIdType(item.PeerType))

		var sessionUnreadCount int64
		if v, ok := retMap[peerId.ToString()]; ok {
			sessionUnreadCount = v
		}

		// createMessage message list
		retList = append(retList, &response.ContactEntity{
			OwnerID:      item.OwnerID,
			OwnerType:    gmodel.ContactIdType(item.OwnerType),
			PeerID:       item.PeerID,
			PeerType:     gmodel.ContactIdType(item.PeerType),
			PeerAck:      gmodel.PeerAckStatus(item.PeerAck),
			VersionID:    item.VersionID,
			SortKey:      item.SortKey,
			Status:       gmodel.ContactStatus(item.Status),
			Labels:       item.Labels,
			LastMsg:      nil,
			UnreadMsgNum: sessionUnreadCount,
		})
	}

	// 返回之前进行重新排序
	//sort.Sort(response.ContactSortBySortKey(retList))
	sort.Slice(retList, func(i, j int) bool {
		return retList[i].VersionID < retList[j].VersionID
	})

	resp.Data = response.FetchContactData{
		ContactList:   retList,
		NextVersionId: maxVersionId,
		HasMore:       len(list) == limit,
	}

	return
}

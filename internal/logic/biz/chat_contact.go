package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"github.com/zhixunjie/im-fun/internal/logic/api"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
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
func (b *ContactUseCase) Fetch(ctx context.Context, req *request.ContactFetchReq) (rsp *response.ContactFetchRsp, err error) {
	rsp = new(response.ContactFetchRsp)
	logHead := fmt.Sprintf("Fetch|req=%v", req)
	limit := 50

	// check params
	err = b.checkParamsFetch(ctx, req)
	if err != nil {
		return
	}

	// get: contact list
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

	lastMsgIds := lo.Map(list, func(item *model.Contact, index int) uint64 {
		return item.LastMsgID
	})
	lastMsgMap, err := b.repoMessage.BatchGetByMsgIds(ctx, lastMsgIds)
	if err != nil {
		logging.Errorf(logHead+"BatchGetByMsgIds err=%v", err)
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
		// build data
		row := &response.ContactEntity{
			OwnerID:      item.OwnerID,
			OwnerType:    gmodel.ContactIdType(item.OwnerType),
			PeerID:       item.PeerID,
			PeerType:     gmodel.ContactIdType(item.PeerType),
			PeerAck:      gmodel.PeerAckStatus(item.PeerAck),
			VersionID:    item.VersionID,
			SortKey:      item.SortKey,
			Status:       gmodel.ContactStatus(item.Status),
			Labels:       item.Labels,
			UnreadMsgNum: sessionUnreadCount,
			LastMsg:      nil,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		}

		// supply last msg info
		if v, ok := lastMsgMap[item.LastMsgID]; ok {
			body := new(format.MsgBody)
			tmpErr := json.Unmarshal([]byte(v.Content), &body)
			if tmpErr == nil {
				row.LastMsg = &response.MsgEntity{
					MsgID:     v.MsgID,
					SeqID:     v.SeqID,
					MsgBody:   body,
					SessionID: v.SessionID,
					SenderID:  v.SenderID,
					SendType:  gmodel.ContactIdType(v.SenderType),
					VersionID: v.VersionID,
					SortKey:   v.SortKey,
					Status:    gmodel.MsgStatus(v.Status),
					HasRead:   gmodel.MsgReadStatus(v.HasRead),
				}
			}
		}

		retList = append(retList, row)
	}

	// 返回之前进行重新排序
	//sort.Sort(response.ContactSortBySortKey(retList))
	sort.Slice(retList, func(i, j int) bool {
		return retList[i].VersionID < retList[j].VersionID
	})

	rsp.Data = &response.FetchContactData{
		ContactList:   retList,
		NextVersionId: maxVersionId,
		HasMore:       len(list) == limit,
	}

	return
}

func (b *ContactUseCase) checkParamsFetch(ctx context.Context, req *request.ContactFetchReq) error {
	// check: sender
	if req.Owner == nil {
		return api.ErrSenderOrReceiverNotAllow
	}
	if req.Owner.Id() == 0 {
		return api.ErrSenderOrReceiverNotAllow
	}

	//allowOwnerType := []gen_id.ContactIdType{
	//	gen_id.TypeUser,
	//}
	//
	//// check: owner type
	//if !lo.Contains(allowOwnerType, req.Owner.Type()) {
	//	return api.ErrSenderTypeNotAllow
	//}

	return nil
}

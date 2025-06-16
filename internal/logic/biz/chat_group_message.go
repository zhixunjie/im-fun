package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/zhixunjie/im-fun/internal/logic/api"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/format"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"github.com/zhixunjie/im-fun/pkg/goredis/distrib_lock"
	k "github.com/zhixunjie/im-fun/pkg/goredis/key"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"github.com/zhixunjie/im-fun/pkg/routine"
	"github.com/zhixunjie/im-fun/pkg/utils"
	"gorm.io/gorm"
	"math"
	"sort"
	"time"
	"unicode/utf8"
)

type GroupMessageUseCase struct {
	repoGroupMessage  *data.GroupMessageRepo
	repoContact       *data.ContactRepo
	repoMessageHelper *data.MessageRepo
}

func NewGroupMessageUseCase(repoGroupMessage *data.GroupMessageRepo, repoContact *data.ContactRepo, repoMessageHelper *data.MessageRepo) *GroupMessageUseCase {
	return &GroupMessageUseCase{
		repoGroupMessage:  repoGroupMessage,
		repoContact:       repoContact,
		repoMessageHelper: repoMessageHelper,
	}
}

// Send 发送消息
func (b *GroupMessageUseCase) Send(ctx context.Context, req *request.GroupMessageSendReq) (rsp *response.GroupMessageSendRsp, err error) {
	rsp = new(response.GroupMessageSendRsp)
	sender := req.Sender
	receiver := req.Receiver
	logHead := fmt.Sprintf("Send|sender=%v,receiver=%v|", sender, receiver)

	// check params
	err = b.checkParamsSend(ctx, req)
	if err != nil {
		return
	}

	// 1. create sender's contact if not exists
	var senderContact, receiverContact *model.ChatContact
	if b.needCreateContact(logHead, sender) {
		if !lo.Contains(req.InvisibleList, cast.ToString(req.Sender.GetId())) {
			senderContact, err = b.repoContact.CreateNotExists(logHead, &model.BuildContactParams{Owner: sender, Peer: receiver})
			if err != nil {
				return
			}
		}
	}
	// 2. create receiver's contact if not exists
	if b.needCreateContact(logHead, receiver) {
		if !lo.Contains(req.InvisibleList, cast.ToString(req.Receiver.GetId())) {
			receiverContact, err = b.repoContact.CreateNotExists(logHead, &model.BuildContactParams{Owner: receiver, Peer: sender})
			if err != nil {
				return
			}
		}
	}
	// 3. create message（无扩散）🔥
	msg, err := b.createMessage(ctx, logHead, req)
	if err != nil {
		return
	}
	currMsgId := msg.MsgID

	routine.Go(ctx, func() {
		// 增加未读数: 先save db，再incr cache，保证尽快执行
		// TODO: 确定群聊场景没有问题
		if !lo.Contains(req.InvisibleList, cast.ToString(req.Receiver.GetId())) {
			_ = b.repoMessageHelper.IncrUnreadAfterSend(ctx, logHead, receiver, sender, 1)
		}
		// update contact's info（写扩散）
		if senderContact != nil {
			err = b.repoContact.UpdateLastMsgId(ctx, logHead, senderContact.ID, sender, currMsgId, gmodel.PeerNotAck)
			if err != nil {
				return
			}
		}
		// update contact's info（写扩散）
		if receiverContact != nil {
			err = b.repoContact.UpdateLastMsgId(ctx, logHead, receiverContact.ID, receiver, currMsgId, gmodel.PeerAcked)
			if err != nil {
				return
			}
			if err != nil {
				return
			}
		}
	})

	rsp.Data = &response.GroupMessageSendData{
		MsgID:       msg.MsgID,
		SeqID:       msg.SeqID,
		VersionID:   msg.VersionID,
		SortKey:     msg.SortKey,
		SessionId:   msg.SessionID,
		UnreadCount: 0,
	}
	return
}

// Fetch 拉取消息列表
func (b *GroupMessageUseCase) Fetch(ctx context.Context, req *request.GroupMessageFetchReq) (rsp *response.GroupMessageFetchRsp, err error) {
	rsp = new(response.GroupMessageFetchRsp)
	logHead := fmt.Sprintf("Fetch|req=%v", req)
	pivotVersionId := req.VersionId
	limit := 50

	// check params
	err = b.checkParamsFetch(ctx, req)
	if err != nil {
		return
	}

	// get: contact info
	owner := req.Owner
	peer := req.Peer
	contactInfo, err := b.repoContact.Info(owner, peer)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = api.ErrContactNotExists
		}
		logging.Error(logHead+"repoContact Info,err=%v", err)
		return
	}

	// get: last msg info
	var lastDelMsgVersionId uint64
	var messageInfo *model.ChatMessage
	if contactInfo.LastDelMsgID > 0 {
		messageInfo, err = b.repoGroupMessage.Info(contactInfo.LastDelMsgID)
		if err != nil {
			return
		}
		lastDelMsgVersionId = messageInfo.VersionID
	}

	// get: message list
	sessionId := gmodel.NewSessionId(owner, peer)
	list, err := b.repoGroupMessage.RangeList(&model.FetchMsgRangeParams{
		FetchType:           req.FetchType,
		SessionId:           sessionId,
		LastDelMsgVersionId: lastDelMsgVersionId,
		PivotVersionId:      pivotVersionId,
		Limit:               limit,
		Owner:               owner,
		Peer:                peer,
	})
	if err != nil {
		return
	}
	if len(list) == 0 {
		return
	}

	// rebuild list
	var retList []*response.GroupMsgEntity
	minVersionId := uint64(math.MaxUint64)
	maxVersionId := uint64(0)
	var tmpList []uint64
	for _, item := range list {
		minVersionId = utils.Min(minVersionId, item.VersionID)
		maxVersionId = utils.Max(maxVersionId, item.VersionID)

		// 过滤：invisible list
		if len(item.InvisibleList) > 0 {
			err = json.Unmarshal([]byte(item.InvisibleList), &tmpList)
			if err != nil {
				continue
			}

			if lo.Contains(tmpList, req.Owner.GetId()) {
				continue
			}
		}

		// build message list
		body := &format.MsgBody{MsgContent: format.NewMsgContent(format.MsgType(item.MsgType))}
		tmpErr := json.Unmarshal([]byte(item.Content), body)
		if tmpErr != nil {
			logging.Error(logHead+"unmarshal msg body fail,err=%v", err)
			continue
		}
		retList = append(retList, &response.GroupMsgEntity{
			MsgID:     item.MsgID,
			SeqID:     item.SeqID,
			MsgBody:   body,
			SessionID: item.SessionID,
			SenderID:  item.SenderID,
			SendType:  gmodel.ContactIdType(item.SenderType),
			VersionID: item.VersionID,
			SortKey:   item.SortKey,
			Status:    gmodel.MsgStatus(item.Status),
			HasRead:   gmodel.MsgReadStatus(item.HasRead),
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	// get: nextVersionId
	var nextVersionId uint64
	switch req.FetchType {
	case gmodel.FetchTypeBackward: // 1. 拉取历史消息
		nextVersionId = minVersionId
	case gmodel.FetchTypeForward: // 2. 拉取最新消息
		nextVersionId = maxVersionId
	default:
		return
	}

	routine.Go(ctx, func() {
		// 减少未读数: 先read db，再decr cache
		// TODO: 确定群聊场景没有问题
		_ = b.repoMessageHelper.DecrUnreadAfterFetch(ctx, logHead, owner, peer, int64(len(retList)))
	})

	// sort: 返回之前进行重新排序
	//sort.Sort(response.GroupMessageSortByVersion(retList))
	sort.Slice(retList, func(i, j int) bool {
		return retList[i].VersionID < retList[j].VersionID
	})

	rsp.Data = &response.GroupMessageFetchData{
		MsgList:       retList,
		NextVersionId: nextVersionId,
		HasMore:       len(list) == limit,
	}

	return
}

// createMessage 构建消息体
func (b *GroupMessageUseCase) createMessage(ctx context.Context, logHead string, req *request.GroupMessageSendReq) (msg *model.ChatMessage, err error) {
	logHead += "createMessage|"
	mem := b.repoMessageHelper.RedisClient
	sender := req.Sender
	receiver := req.Receiver

	// exchange：InvisibleList
	var bInvisibleList []byte
	if len(req.InvisibleList) > 0 {
		bInvisibleList, err = json.Marshal(req.InvisibleList)
		if err != nil {
			logging.Errorf(logHead+"Marshal error=%v", err)
			return
		}
	}

	// exchange：MsgContent
	content, err := json.Marshal(req.MsgBody)
	if err != nil {
		logging.Errorf(logHead+"Marshal error=%v", err)
		return
	}

	// message: gen session id
	sessionId := gmodel.NewSessionId(sender, receiver)

	// note: 同一个消息timeline的版本变动，需要加锁
	// 保证数据库记录中的 msg_id 与 session_id 是递增的
	lockKey := data.TimelineMessageLock.Format(k.M{"session_id": sessionId})
	redisSpinLock := distrib_lock.NewSpinLock(mem, lockKey, 5*time.Second, &distrib_lock.SpinOption{Interval: 50 * time.Millisecond, Times: 40})
	if err = redisSpinLock.AcquireWithTimes(); err != nil {
		logging.Errorf(logHead+"acquire fail,lockKey=%v,err=%v", lockKey, err)
		return
	}
	defer func() { _ = redisSpinLock.Release() }()
	logging.Infof(logHead+"acquire success,lockKey=%v", lockKey)

	// a) generate message's msg_id
	msgId, err := gen_id.NewMsgId(ctx, &gen_id.MsgIdParams{
		Mem: mem,
		Id1: sender,
		Id2: receiver,
	})
	if err != nil {
		logging.Errorf(logHead+"gen MsgID error=%v", err)
		return
	}

	// b) generate message's version_id
	versionId, err := gen_id.NewMsgVersionId(ctx, &gen_id.MsgVerParams{
		Mem: mem,
		Id1: sender,
		Id2: receiver,
	})
	if err != nil {
		logging.Errorf(logHead+"gen VersionID error=%v", err)
		return
	}

	// build message
	msg = &model.ChatMessage{
		MsgID:         msgId,                          // 唯一id（服务端）
		SeqID:         req.SeqId,                      // 唯一id（客户端）
		MsgType:       uint32(req.MsgBody.MsgType),    // 消息类型
		Content:       string(content),                // 消息内容
		SessionID:     string(sessionId),              // 会话ID
		SenderID:      req.Sender.GetId(),             // 发送者ID
		SenderType:    uint32(req.Sender.GetType()),   // 发送者的用户类型
		VersionID:     versionId,                      // 版本ID
		SortKey:       versionId,                      // sort_key的值等同于version_id
		Status:        uint32(gmodel.MsgStatusNormal), // 状态正常
		HasRead:       uint32(gmodel.MsgRead),         // TODO: 已读（功能还没做好）
		InvisibleList: string(bInvisibleList),         // 不可见的列表
	}
	err = b.repoGroupMessage.Create(logHead, msg)
	if err != nil {
		return
	}

	return
}

// 双方通信时，判断是否需要创建对方的Contact
func (b *GroupMessageUseCase) needCreateContact(logHead string, id *gmodel.ComponentId) bool {
	logHead += fmt.Sprintf("needCreateContact|id=%v|", id)

	noNeedCreate := []gmodel.ContactIdType{
		gmodel.TypeRobot,
		gmodel.TypeSystem,
	}

	// 如果用户是指定类型，那么不需要创建他的contact信息（比如：机器人）
	if lo.Contains(noNeedCreate, id.GetType()) || id.IsGroup() {
		logging.Info(logHead + "do not need create contact")
		return false
	}

	return true
}

// 限制：发送者和接受者的类型
func (b *GroupMessageUseCase) checkParamsSend(ctx context.Context, req *request.GroupMessageSendReq) error {
	// check: user
	if req.Sender == nil || req.Receiver == nil {
		return api.ErrSenderOrReceiverNotAllow
	}
	if req.Sender.GetId() == 0 || req.Receiver.GetId() == 0 {
		return api.ErrSenderOrReceiverNotAllow
	}
	if req.Sender.Equal(req.Receiver) {
		return fmt.Errorf("ID equal %w", api.ErrSenderOrReceiverNotAllow)
	}
	if !req.Sender.IsGroup() && !req.Receiver.IsGroup() {
		return fmt.Errorf("group not allowed %w", api.ErrSenderOrReceiverNotAllow)
	}

	//allowSenderType := []gen_id.ContactIdType{
	//	gen_id.TypeUser,
	//	gen_id.TypeRobot,
	//	gen_id.TypeSystem,
	//}
	//
	//allowReceiverType := []gen_id.ContactIdType{
	//	gen_id.TypeUser,
	//	gen_id.TypeRobot,
	//	gen_id.TypeGroup,
	//}
	//
	//// check: sender type
	//if !lo.Contains(allowSenderType, req.Sender.GetType()) {
	//	return api.ErrSenderTypeNotAllow
	//}
	//// check: receiver type
	//if !lo.Contains(allowReceiverType, req.Receiver.GetType()) {
	//	return api.ErrReceiverTypeNotAllow
	//}

	// check: message body
	if req.MsgBody == nil {
		return api.ErrMessageBodyNotAllow
	}

	// check: message length
	content, err := json.Marshal(req.MsgBody)
	if utf8.RuneCount(content) > 2048 {
		return fmt.Errorf("%v(content is too long)", api.ErrMessageContentNotAllowed)
	}

	// check: message type
	typeLimit := []format.MsgType{
		format.MsgTypeCustom,
		format.MsgTypeText,
		format.MsgTypeImage,
		format.MsgTypeVideo,
		format.MsgTypeTips,
	}
	if !lo.Contains(typeLimit, req.MsgBody.MsgType) {
		return api.ErrMessageTypeNotAllowed
	}

	// check: message content
	msgContent, err := format.DecodeMsgBody(req.MsgBody)
	if err != nil {
		return api.ErrMessageBodyDecodedFailed
	}
	switch v := msgContent.(type) {
	case *format.CustomContent:
		if v.Data == "" {
			return fmt.Errorf("%v(text is empty)", api.ErrMessageContentNotAllowed)
		}
	case *format.TextContent:
		if v.Text == "" {
			return fmt.Errorf("%v(text is empty)", api.ErrMessageContentNotAllowed)
		}
	case *format.ImageContent:
		if len(v.ImageInfos) == 0 {
			return fmt.Errorf("%v(image array is empty)", api.ErrMessageContentNotAllowed)
		}
	case *format.VideoContent:
		if v.VideoUrl == "" {
			return fmt.Errorf("%v(video url is empty)", api.ErrMessageContentNotAllowed)
		}
		if v.VideoSecond == 0 {
			return fmt.Errorf("%v(video second is zero)", api.ErrMessageContentNotAllowed)
		}
	case *format.TipsContent:
		if v.Text == "" {
			return fmt.Errorf("%v(tip's text is empty)", api.ErrMessageContentNotAllowed)
		}
	}

	// TODO: 频率控制、敏感词控制

	return nil
}

func (b *GroupMessageUseCase) checkParamsFetch(ctx context.Context, req *request.GroupMessageFetchReq) error {
	// check: user
	if req.Owner == nil || req.Peer == nil {
		return api.ErrSenderOrReceiverNotAllow
	}
	if req.Owner.GetId() == 0 || req.Peer.GetId() == 0 {
		return api.ErrSenderOrReceiverNotAllow
	}
	if req.Owner.Equal(req.Peer) {
		return fmt.Errorf("ID equal %w", api.ErrSenderOrReceiverNotAllow)
	}
	if !req.Owner.IsGroup() && !req.Peer.IsGroup() {
		return fmt.Errorf("group not allowed %w", api.ErrSenderOrReceiverNotAllow)
	}

	//allowOwnerType := []gen_id.ContactIdType{
	//	gen_id.TypeUser,
	//}
	//
	//allowPeerType := []gen_id.ContactIdType{
	//	gen_id.TypeUser,
	//	gen_id.TypeRobot,
	//	gen_id.TypeGroup,
	//}
	//
	//// check: owner type
	//if !lo.Contains(allowOwnerType, req.Owner.GetType()) {
	//	return api.ErrSenderTypeNotAllow
	//}
	//// check: peer type
	//if !lo.Contains(allowPeerType, req.Peer.GetType()) {
	//	return api.ErrReceiverTypeNotAllow
	//}

	return nil
}

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
	"github.com/zhixunjie/im-fun/internal/logic/data/cache"
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

type MessageUseCase struct {
	repoMessage *data.MessageRepo
	repoContact *data.ContactRepo
}

func NewMessageUseCase(repoMessage *data.MessageRepo, repoContact *data.ContactRepo) *MessageUseCase {
	return &MessageUseCase{
		repoMessage: repoMessage,
		repoContact: repoContact,
	}
}

// SendSimpleCustomMessage 简化接口：发送自定义消息
func (b *MessageUseCase) SendSimpleCustomMessage(ctx context.Context, sender, receiver *gmodel.ComponentId, d string) (rsp *response.MessageSendRsp, err error) {
	return b.Send(ctx, &request.MessageSendReq{
		SeqId:    uint64(gmodel.NewSeqId()),
		Sender:   sender,
		Receiver: receiver,
		MsgBody: &format.MsgBody{
			MsgType: format.MsgTypeCustom,
			MsgContent: &format.CustomContent{
				Data: d,
			},
		},
	})
}

// Send 发送消息
func (b *MessageUseCase) Send(ctx context.Context, req *request.MessageSendReq) (rsp *response.MessageSendRsp, err error) {
	rsp = new(response.MessageSendRsp)
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
		if !lo.Contains(req.InvisibleList, cast.ToString(req.Sender.Id())) {
			senderContact, err = b.repoContact.CreateNotExists(logHead, &model.BuildContactParams{Owner: sender, Peer: receiver})
			if err != nil {
				return
			}
		}
	}
	// 2. create receiver's contact if not exists
	if b.needCreateContact(logHead, receiver) {
		if !lo.Contains(req.InvisibleList, cast.ToString(req.Receiver.Id())) {
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
		if !lo.Contains(req.InvisibleList, cast.ToString(req.Receiver.Id())) {
			_ = b.repoMessage.IncrUnreadAfterSend(ctx, logHead, receiver, sender, 1)
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

	rsp.Data = &response.SendMsgRespData{
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
func (b *MessageUseCase) Fetch(ctx context.Context, req *request.MessageFetchReq) (rsp *response.MessageFetchRsp, err error) {
	rsp = new(response.MessageFetchRsp)
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
		messageInfo, err = b.repoMessage.Info(contactInfo.LastDelMsgID)
		if err != nil {
			return
		}
		lastDelMsgVersionId = messageInfo.VersionID
	}

	// get: message list
	sessionId := gmodel.NewSessionId(owner, peer)
	list, err := b.repoMessage.RangeList(&model.FetchMsgRangeParams{
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
	var retList []*response.MsgEntity
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

			if lo.Contains(tmpList, req.Owner.Id()) {
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
		retList = append(retList, &response.MsgEntity{
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
		_ = b.repoMessage.DecrUnreadAfterFetch(ctx, logHead, owner, peer, int64(len(retList)))
	})

	// sort: 返回之前进行重新排序
	//sort.Sort(response.MessageSortByVersion(retList))
	sort.Slice(retList, func(i, j int) bool {
		return retList[i].VersionID < retList[j].VersionID
	})

	rsp.Data = &response.FetchMsgData{
		MsgList:       retList,
		NextVersionId: nextVersionId,
		HasMore:       len(list) == limit,
	}

	return
}

// Recall 某条消息撤回（两边的聊天记录都需要撤回）
// 核心：更新status和version_id（需要通知对方，所以需要更新version_id）
func (b *MessageUseCase) Recall(ctx context.Context, req *request.MessageRecallReq) (rsp *response.MessageRecallRsp, err error) {
	rsp = new(response.MessageRecallRsp)
	logHead := "Recall|"
	senderId := req.Sender

	// invoke common method
	err = b.updateMsgIdStatus(ctx, logHead, req.MsgID, gmodel.MsgStatusRecall, senderId)
	if err != nil {
		return
	}
	return
}

// DelBothSide 删除一条消息（两边的聊天记录都需要删除）
// 核心：更新status和version_id（需要通知对方，所以需要更新version_id）
func (b *MessageUseCase) DelBothSide(ctx context.Context, req *request.MessageDelBothSideReq) (rsp *response.DelBothSideRsp, err error) {
	rsp = new(response.DelBothSideRsp)
	logHead := "DelBothSide|"
	senderId := req.Sender

	// invoke common method
	err = b.updateMsgIdStatus(ctx, logHead, req.MsgID, gmodel.MsgStatusDeleted, senderId)
	if err != nil {
		return
	}
	return
}

// DelOneSide 删除一条消息（只有一边的聊天记录是不可见，另外一边可见）
// 核心：更新 invisible_list
func (b *MessageUseCase) DelOneSide(ctx context.Context, req *request.MessageDelOneSideReq) (rsp *response.DelOneSideRsp, err error) {
	rsp = new(response.DelOneSideRsp)
	return
}

func (b *MessageUseCase) checkParamsClearHistory(ctx context.Context, req *request.MessageClearHistoryReq) (lastDelMsgId uint64, err error) {
	owner := req.Owner
	peer := req.Peer
	lastDelMsgId = req.MsgID

	// check: user
	if req.Owner == nil || req.Peer == nil {
		err = api.ErrSenderOrReceiverNotAllow
		return
	}
	if req.Owner.Id() == 0 || req.Peer.Id() == 0 {
		err = api.ErrSenderOrReceiverNotAllow
		return
	}
	if req.Owner.Equal(req.Peer) {
		err = fmt.Errorf("ID equal %w", api.ErrSenderOrReceiverNotAllow)
		return
	}

	// 此接口不适合群聊场景
	if owner.IsGroup() || peer.IsGroup() {
		err = errors.New("clear history not allow group")
		return
	}

	// 聊天记录的清空：
	// 1. 如果指定 msgId，则对该 msgId 为界限进行清空；
	// 2. 如果没有指定，则通过 contact 记录的最后一条消息进行清空（通信双方都有各自的 last_msg_id 和 last_del_msg_id）；
	if lastDelMsgId > 0 {
		var msgInfo *model.ChatMessage
		msgInfo, err = b.repoMessage.Info(lastDelMsgId)
		if err != nil {
			err = fmt.Errorf("get message info error: %w", err)
			return
		}
		// check: 是否属于通信双方的信息
		var receiver *gmodel.ComponentId
		receiver, err = gmodel.SessionId(msgInfo.SessionID).ParsePeerId(owner)
		if err != nil {
			err = errors.New("parse session id error")
			return
		}
		if !receiver.Equal(peer) {
			err = errors.New("lastDelMsgId not match peer")
			return
		}
	} else {
		// get: contact info
		var contactInfo *model.ChatContact
		contactInfo, err = b.repoContact.Info(owner, peer)
		if err != nil {
			err = fmt.Errorf("get contact info error: %w", err)
			return
		}
		lastDelMsgId = contactInfo.LastMsgID
	}

	return
}

// ClearHistory
// 清空聊天记录（批量清空），比如：删除联系人后，需要把自己这边的聊天记录给清空（如果需要清空双方，需要把双方的 contact 都做标记）
// 核心：更新Contact的 last_del_msg_id 为 last_msg_id
func (b *MessageUseCase) ClearHistory(ctx context.Context, req *request.MessageClearHistoryReq) (rsp *response.MessageClearHistoryRsp, err error) {
	rsp = new(response.MessageClearHistoryRsp)
	logHead := fmt.Sprintf("ClearHistory|")
	owner := req.Owner
	peer := req.Peer
	mem := b.repoMessage.RedisClient

	// check params
	lastDelMsgId, err := b.checkParamsClearHistory(ctx, req)
	if err != nil {
		logging.Error(logHead+"checkParamsClearHistory error", err)
		return
	}

	// contact: gen version_id
	versionId, err := gen_id.NewContactVersionId(ctx, &gen_id.ContactVerParams{
		Mem:   mem,
		Owner: owner,
	})
	if err != nil {
		logging.Errorf(logHead+"gen VersionID error=%v", err)
		return
	}

	// update to db
	err = b.repoContact.UpdateLastDelMsg(lastDelMsgId, versionId, owner, peer)
	if err != nil {
		logging.Errorf(logHead+"UpdateLastDelMsg error=%v", err)
		return
	}
	rsp.LastDelMsgID = lastDelMsgId

	return
}

// createMessage 构建消息体
func (b *MessageUseCase) createMessage(ctx context.Context, logHead string, req *request.MessageSendReq) (msg *model.ChatMessage, err error) {
	logHead += "createMessage|"
	mem := b.repoMessage.RedisClient
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
	lockKey := cache.TimelineMessageLock.Format(k.M{"session_id": sessionId})
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
		SenderID:      req.Sender.Id(),                // 发送者ID
		SenderType:    uint32(req.Sender.Type()),      // 发送者的用户类型
		VersionID:     versionId,                      // 版本ID
		SortKey:       versionId,                      // sort_key的值等同于version_id
		Status:        uint32(gmodel.MsgStatusNormal), // 状态正常
		HasRead:       uint32(gmodel.MsgRead),         // TODO: 已读（功能还没做好）
		InvisibleList: string(bInvisibleList),         // 不可见的列表
	}
	err = b.repoMessage.Create(logHead, msg)
	if err != nil {
		return
	}

	return
}

// 双方通信时，判断是否需要创建对方的Contact
func (b *MessageUseCase) needCreateContact(logHead string, id *gmodel.ComponentId) bool {
	logHead += fmt.Sprintf("needCreateContact|id=%v|", id)

	noNeedCreate := []gmodel.ContactIdType{
		gmodel.TypeRobot,
		gmodel.TypeSystem,
	}

	// 如果用户是指定类型，那么不需要创建他的contact信息（比如：机器人）
	if lo.Contains(noNeedCreate, id.Type()) || id.IsGroup() {
		logging.Info(logHead + "do not need create contact")
		return false
	}

	return true
}

// 限制：发送者和接受者的类型
func (b *MessageUseCase) checkParamsSend(ctx context.Context, req *request.MessageSendReq) error {
	// check: user
	if req.Sender == nil || req.Receiver == nil {
		return api.ErrSenderOrReceiverNotAllow
	}
	if req.Sender.Id() == 0 || req.Receiver.Id() == 0 {
		return api.ErrSenderOrReceiverNotAllow
	}
	if req.Sender.Equal(req.Receiver) {
		return fmt.Errorf("ID equal %w", api.ErrSenderOrReceiverNotAllow)
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
	//if !lo.Contains(allowSenderType, req.Sender.Type()) {
	//	return api.ErrSenderTypeNotAllow
	//}
	//// check: receiver type
	//if !lo.Contains(allowReceiverType, req.Receiver.Type()) {
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

func (b *MessageUseCase) checkParamsFetch(ctx context.Context, req *request.MessageFetchReq) error {
	// check: user
	if req.Owner == nil || req.Peer == nil {
		return api.ErrSenderOrReceiverNotAllow
	}
	if req.Owner.Id() == 0 || req.Peer.Id() == 0 {
		return api.ErrSenderOrReceiverNotAllow
	}
	if req.Owner.Equal(req.Peer) {
		return fmt.Errorf("ID equal %w", api.ErrSenderOrReceiverNotAllow)
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
	//if !lo.Contains(allowOwnerType, req.Owner.Type()) {
	//	return api.ErrSenderTypeNotAllow
	//}
	//// check: peer type
	//if !lo.Contains(allowPeerType, req.Peer.Type()) {
	//	return api.ErrReceiverTypeNotAllow
	//}

	return nil
}

// updateMsgIdStatus 通用的方法，用于更新消息的状态和版本ID
func (b *MessageUseCase) updateMsgIdStatus(ctx context.Context, logHead string, msgId model.BigIntType, status gmodel.MsgStatus, sender *gmodel.ComponentId) (err error) {
	logHead += fmt.Sprintf("updateMsgIdStatus,msgId=%v,status=%v|", msgId, status)

	// get: message
	msgInfo, err := b.repoMessage.Info(msgId)
	if err != nil {
		logging.Errorf(logHead+"repoMessage Info error=%v", err)
		return
	}
	sessionId := msgInfo.SessionID

	// save to db
	fn := func(versionId uint64) (err error) {
		// update to db
		err = b.repoMessage.UpdateMsgVerAndStatus(logHead, msgId, versionId, status)
		if err != nil {
			return
		}
		return
	}
	// 加锁生成 version_id，然后执行回调函数
	err = b.updateMsgVersion(ctx, logHead, sessionId, sender, fn)
	if err != nil {
		return
	}

	return
}

type Fn func(versionId uint64) (err error)

// updateMsgVersion 加锁 => 生成 version_id => 执行回调函数
func (b *MessageUseCase) updateMsgVersion(ctx context.Context, logHead string, sessionId string, sender *gmodel.ComponentId, fn Fn) (err error) {
	mem := b.repoMessage.RedisClient

	// 提取: 接收者的信息
	receiver, err := gmodel.SessionId(sessionId).ParsePeerId(sender)
	if err != nil {
		err = fmt.Errorf("ParsePeerId failed: %w", err)
		return
	}

	// note: 同一个消息timeline的版本变动，需要加锁
	lockKey := cache.TimelineMessageLock.Format(k.M{"session_id": sessionId})
	redisSpinLock := distrib_lock.NewSpinLock(mem, lockKey, 5*time.Second, &distrib_lock.SpinOption{Interval: 50 * time.Millisecond, Times: 40})
	if err = redisSpinLock.AcquireWithTimes(); err != nil {
		logging.Errorf(logHead+"acquire fail,lockKey=%v,err=%v", lockKey, err)
		return
	}
	defer func() { _ = redisSpinLock.Release() }()
	logging.Infof(logHead+"acquire success,lockKey=%v", lockKey)

	// generate message's version_id
	versionId, err := gen_id.NewMsgVersionId(ctx, &gen_id.MsgVerParams{
		Mem: mem,
		Id1: sender,
		Id2: receiver,
	})
	if err != nil {
		logging.Errorf(logHead+"gen VersionID error=%v", err)
		return
	}

	return fn(versionId)
}

package service

import (
	"context"
	"github.com/zhixunjie/im-fun/internal/logic/dao"
	"github.com/zhixunjie/im-fun/internal/logic/model"
	"github.com/zhixunjie/im-fun/internal/logic/model/request"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"time"
)

func (svc *Service) Ping() {

}

// 消息发送方的会话
func transformSenderContact(ctx context.Context, req *request.SendMsgReq, currTimestamp int64, msgId uint64) (model.Contact, error) {
	var defaultRet model.Contact

	// get version_id（区别的地方）
	versionId, err := gen_id.GetContactVersionId(ctx, currTimestamp, req.SendId)
	if err != nil {
		return defaultRet, err
	}

	// query contact（区别的地方，获取"消息发送方"的会话）
	contact, err := dao.QueryContactById(req.SendId, req.ReceiveId)
	if err != nil {
		return defaultRet, err
	}
	// 新增：需要执行的逻辑
	if contact.Id == 0 {
		contact.PeerType = model.PeerNotExist
		contact.PeerAck = model.PeerNotAck
		contact.CreatedAt = time.Now()
	}
	// 新增 or 更新：都要执行的逻辑
	contact.OwnerId = req.SendId   // 会话的所有者
	contact.PeerId = req.ReceiveId // 会话的对方
	contact.LastMsgId = msgId      // 双方聊天记录中，最新一次发送的消息id
	contact.VersionId = versionId  // 版本号（用于拉取会话框）
	contact.SortKey = versionId    // sort_key的值等同于version_id
	contact.PeerType = req.SendType
	contact.Status = model.ContactStatusNormal
	contact.UpdatedAt = time.Now()

	return contact, nil
}

// 消息接收方的会话
func transformReceiverContact(ctx context.Context, req *request.SendMsgReq, currTimestamp int64, msgId uint64) (model.Contact, error) {
	var defaultRet model.Contact

	// get version_id（区别的地方）
	versionId, err := gen_id.GetContactVersionId(ctx, currTimestamp, req.ReceiveId)
	if err != nil {
		return defaultRet, err
	}

	// query contact（区别的地方，获取"消息接收方"的会话）
	contact, err := dao.QueryContactById(req.ReceiveId, req.SendId)
	if err != nil {
		return defaultRet, err
	}
	// 新增：需要执行的逻辑（区别的地方）
	if contact.Id == 0 {
		contact.PeerType = model.PeerNotExist
		contact.PeerAck = model.PeerAck
		contact.CreatedAt = time.Now()
	}
	// 新增 or 更新：都要执行的逻辑（区别的地方）
	contact.OwnerId = req.ReceiveId // 会话的所有者
	contact.PeerId = req.SendId     // 会话的对方
	contact.LastMsgId = msgId       // 双方聊天记录中，最新一次发送的消息id
	contact.VersionId = versionId   // 版本号（用于拉取会话框）
	contact.SortKey = versionId     // sort_key的值等同于version_id
	contact.PeerType = req.ReceiveType
	contact.Status = model.ContactStatusNormal
	contact.UpdatedAt = time.Now()

	return contact, nil
}

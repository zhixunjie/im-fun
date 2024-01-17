package biz

import (
	"context"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"gorm.io/gorm"
)

type ContactUseCase struct {
	repo *data.ContactRepo
}

func NewContactUseCase(repo *data.ContactRepo) *ContactUseCase {
	return &ContactUseCase{repo: repo}
}

// TransformSender 消息发送方的会话
func (bz *ContactUseCase) TransformSender(ctx context.Context, req *request.SendMsgReq, currTimestamp int64, msgId uint64) (contact *model.Contact, err error) {
	// get version_id（区别的地方）
	versionId, err := gen_id.ContactVersionId(ctx, bz.repo.RedisClient, currTimestamp, req.SendId)
	if err != nil {
		return
	}

	// query contact（区别的地方，获取"消息发送方"的会话）
	contact, err = bz.repo.QueryContactById(req.SendId, req.PeerId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	// 新增：需要执行的逻辑
	if contact.ID == 0 {
		contact.PeerType = model.PeerNotExist
		contact.PeerAck = model.PeerNotAck
	}
	// 新增 or 更新：都要执行的逻辑
	contact.OwnerID = req.SendId  // 会话的所有者
	contact.PeerID = req.PeerId   // 会话的对方
	contact.LastMsgID = msgId     // 双方聊天记录中，最新一次发送的消息id
	contact.VersionID = versionId // 版本号（用于拉取会话框）
	contact.SortKey = versionId   // sort_key的值等同于version_id
	contact.PeerType = req.SenderType
	contact.Status = model.ContactStatusNormal

	return contact, nil
}

// TransformPeer 消息接收方的会话
func (bz *ContactUseCase) TransformPeer(ctx context.Context, req *request.SendMsgReq, currTimestamp int64, msgId uint64) (contact *model.Contact, err error) {
	// get version_id（区别的地方）
	versionId, err := gen_id.ContactVersionId(ctx, bz.repo.RedisClient, currTimestamp, req.PeerId)
	if err != nil {
		return
	}

	// query contact（区别的地方，获取"消息接收方"的会话）
	contact, err = bz.repo.QueryContactById(req.PeerId, req.SendId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	// 新增：需要执行的逻辑（区别的地方）
	if contact.ID == 0 {
		contact.PeerType = model.PeerNotExist
		contact.PeerAck = model.PeerAck
	}
	// 新增 or 更新：都要执行的逻辑（区别的地方）
	contact.OwnerID = req.PeerId  // 会话的所有者
	contact.PeerID = req.SendId   // 会话的对方
	contact.LastMsgID = msgId     // 双方聊天记录中，最新一次发送的消息id
	contact.VersionID = versionId // 版本号（用于拉取会话框）
	contact.SortKey = versionId   // sort_key的值等同于version_id
	contact.PeerType = req.PeerType
	contact.Status = model.ContactStatusNormal

	return
}

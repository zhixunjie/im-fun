package biz

import (
	"context"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"gorm.io/gorm"
)

type ContactUseCase struct {
	repo *data.ContactRepo
}

func NewContactUseCase(repo *data.ContactRepo) *ContactUseCase {
	return &ContactUseCase{repo: repo}
}

func (contactUseCase *ContactUseCase) Transform(ctx context.Context, ownerId, peerId uint64, peerType int32, currTimestamp int64, msgId uint64) (contact *model.Contact, err error) {
	// get version_id
	versionId, err := gen_id.ContactVersionId(ctx, contactUseCase.repo.RedisClient, currTimestamp, ownerId)
	if err != nil {
		return
	}

	// query contact
	contact, err = contactUseCase.repo.QueryContactById(ownerId, peerId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	// 记录不存在：需要创建contact
	if err == gorm.ErrRecordNotFound {
		contact = new(model.Contact)
		contact.PeerType = model.PeerNotExist
		contact.PeerAck = model.PeerAck
	}
	// 新增 or 更新：都要执行的逻辑（区别的地方）
	contact.OwnerID = ownerId     // 会话的所有者
	contact.PeerID = peerId       // 联系人
	contact.PeerType = peerType   // 联系人类型
	contact.LastMsgID = msgId     // 双方聊天记录中，最新一次发送的消息id
	contact.VersionID = versionId // 版本号（用于拉取会话框）
	contact.SortKey = versionId   // sort_key的值等同于version_id
	contact.Status = model.ContactStatusNormal

	return
}

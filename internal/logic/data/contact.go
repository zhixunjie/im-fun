package data

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/query"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"gorm.io/gorm"
)

type ContactRepo struct {
	*Data
}

func NewContactRepo(data *Data) *ContactRepo {
	return &ContactRepo{
		Data: data,
	}
}

func (repo *ContactRepo) TableName(ownerId uint64) (dbName string, tbName string) {
	//// 临时写死
	//if true {
	//	return "", "contact_0"
	//}
	// 分表规则：
	// - 数据库前缀：message_xxx，规则：owner_id 倒数第三位数字就是分库值
	// - 数据表前缀：contact_xxx，规则：owner_id 的最后两位就是分表值
	dbName = fmt.Sprintf("messsage_%v", ownerId%1000/100)
	tbName = fmt.Sprintf("contact_%v", ownerId%model.TotalTableContact)

	return dbName, tbName
}

// InfoWithCache 查询某个会话的信息
func (repo *ContactRepo) InfoWithCache(ownerId uint64, peerId uint64) (*model.Contact, error) {
	// todo 先从cache拿，拿不到再从DB拿

	return repo.Info(ownerId, peerId)
}

// Info 查询某个会话的信息
func (repo *ContactRepo) Info(ownerId uint64, peerId uint64) (row *model.Contact, err error) {
	_, tbName := repo.TableName(ownerId)
	qModel := repo.Db.Contact.Table(tbName)

	row, err = qModel.Where(qModel.OwnerID.Eq(ownerId), qModel.PeerID.Eq(peerId)).Take()
	if err != nil {
		return
	}
	return
}

// Edit 插入/更新记录
func (repo *ContactRepo) Edit(tx *query.Query, row *model.Contact) (err error) {
	_, tbName := repo.TableName(row.OwnerID)
	qModel := tx.Contact.Table(tbName)

	if row.ID == 0 {
		err = qModel.Create(row)
		if err != nil {
			return
		}
	} else {
		_, err = qModel.Updates(row)
		if err != nil {
			return
		}
	}

	return
}

func (repo *ContactRepo) Build(ctx context.Context, params *model.BuildContactParams) (contact *model.Contact, err error) {
	// get version_id
	versionId, err := gen_id.ContactVersionId(ctx, repo.RedisClient, params.OwnerId)
	if err != nil {
		return
	}

	// query contact
	contact, err = repo.Info(params.OwnerId, params.PeerId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	// 记录不存在：需要创建contact
	if err == gorm.ErrRecordNotFound {
		contact = &model.Contact{
			OwnerID: params.OwnerId,
			PeerID:  params.PeerId,
			//PeerType:  params.PeerType, // 注意：暂时不适用请求参数过来的PeerType（适合于：logic -> base的场景）
			PeerType:  int32(model.PeerTypeNormal),
			PeerAck:   params.PeerAck,
			LastMsgID: params.MsgId,
			VersionID: versionId,
			SortKey:   versionId,
			Status:    model.ContactStatusNormal,
		}
	} else {
		contact.LastMsgID = params.MsgId // 双方聊天记录中，最新一次发送的消息id
		contact.VersionID = versionId    // 版本号（用于拉取会话框）
		contact.SortKey = versionId      // sort_key的值等同于version_id
	}
	return
}

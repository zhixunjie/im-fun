package data

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/cache"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/query"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/goredis/distrib_lock"
	k "github.com/zhixunjie/im-fun/pkg/goredis/key"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"gorm.io/gen"
	"gorm.io/gorm"
	"math"
	"time"
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
	// 临时写死
	if true {
		return "", "contact"
	}
	// 分表规则：
	// - 数据库前缀：message_xxx，规则：owner_id 倒数第三位数字就是分库值
	// - 数据表前缀：contact_xxx，规则：owner_id 的最后两位就是分表值
	dbName = fmt.Sprintf("messsage_%v", ownerId%1000/100)
	tbName = fmt.Sprintf("contact_%v", ownerId%model.TotalTableContact)

	return dbName, tbName
}

// InfoWithCache 查询某个会话的信息
func (repo *ContactRepo) InfoWithCache(ownerId *gen_id.ComponentId, peerId *gen_id.ComponentId) (*model.Contact, error) {
	// todo 先从cache拿，拿不到再从DB拿

	return repo.Info("InfoWithCache|", ownerId, peerId)
}

// Info 查询某个会话的信息
func (repo *ContactRepo) Info(logHead string, ownerId *gen_id.ComponentId, peerId *gen_id.ComponentId) (row *model.Contact, err error) {
	_, tbName := repo.TableName(ownerId.Id())
	qModel := repo.Db.Contact.Table(tbName)

	row, err = qModel.Where(
		qModel.OwnerID.Eq(ownerId.Id()),
		qModel.OwnerType.Eq(ownerId.Type()),
		qModel.PeerID.Eq(peerId.Id()),
		qModel.PeerType.Eq(peerId.Type()),
	).Take()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logging.Errorf(logHead+"Take err=%v", err)
		}
		return
	}
	return
}

// Edit 插入/更新记录
func (repo *ContactRepo) Edit(logHead string, tx *query.Query, row *model.Contact) (err error) {
	logHead += fmt.Sprintf("Edit,row=%v|", row)

	_, tbName := repo.TableName(row.OwnerID)
	qModel := tx.Contact.Table(tbName)

	// insert or update ?
	if row.ID == 0 {
		err = qModel.Create(row)
		if err != nil {
			logging.Errorf(logHead+"Create fail,err=%v", err)
			return
		}
		logging.Infof(logHead + "Create success")
	} else {
		var res gen.ResultInfo
		res, err = qModel.Where(qModel.ID.Eq(row.ID)).Limit(1).Updates(row)
		if err != nil {
			logging.Errorf(logHead+"Updates fail,err=%v", err)
			return
		}
		logging.Infof(logHead+"Updates success,RowsAffected=%v", res.RowsAffected)
	}

	return
}

// CreateNotExists 创建会话
func (repo *ContactRepo) CreateNotExists(logHead string, params *model.BuildContactParams) (contact *model.Contact, err error) {
	logHead += "CreateNotExists|"
	_, tbName := repo.TableName(params.OwnerId.Id())
	qModel := repo.Db.Contact.Table(tbName)

	// query contact
	contact, err = repo.Info(logHead, params.OwnerId, params.PeerId)
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Errorf(logHead+"Info error=%v", err)
		return
	}

	// insert if not exists
	if err == gorm.ErrRecordNotFound {
		contact = &model.Contact{
			OwnerID: params.OwnerId.Id(), OwnerType: params.OwnerId.Type(),
			PeerID: params.PeerId.Id(), PeerType: params.PeerId.Type(),
			PeerAck: uint32(params.PeerAck),
			Status:  model.ContactStatusNormal,
		}

		// save to db
		err = qModel.Create(contact)
		if err != nil {
			logging.Errorf(logHead+"Create fail,err=%v,contact=%v", err, contact)
			return
		}
		logging.Infof(logHead+"Create success,contact=%v", contact)
	}
	return
}

// UpdateLastMsgId 更新contact的最后一条消息（发消息）
func (repo *ContactRepo) UpdateLastMsgId(ctx context.Context, logHead string, contactId uint64, lastMsgId uint64, ownerId *gen_id.ComponentId) (err error) {
	logHead += "UpdateLastMsgId|"
	mem := repo.RedisClient
	_, tbName := repo.TableName(ownerId.Id())
	qModel := repo.Db.Contact.Table(tbName)

	// note: 同一用户的会话timeline的版本变动，需要加锁
	lockKey := cache.TimelineContactLock.Format(k.M{"uid": ownerId.Id()})
	redisSpinLock := distrib_lock.NewSpinLock(mem, lockKey, 5*time.Second, &distrib_lock.SpinOption{Interval: 20 * time.Millisecond, Times: 20})
	if err = redisSpinLock.AcquireWithTimes(); err != nil {
		logging.Errorf(logHead+"acquire fail,lockKey=%v,err=%v", lockKey, err)
		return
	}
	defer func() { _ = redisSpinLock.Release() }()
	logging.Infof(logHead+"acquire success,lockKey=%v", lockKey)

	// contact: get version_id
	var versionId uint64
	versionId, err = gen_id.VersionId(ctx, &gen_id.GenVersionParams{
		Mem:            mem,
		GenVersionType: gen_id.GenVersionTypeContact,
		OwnerId:        ownerId,
	})
	if err != nil {
		logging.Errorf(logHead+"gen VersionId error=%v", err)
		return
	}

	// 只更新一部分的字段
	row := &model.Contact{
		LastMsgID: lastMsgId, // 1. 双方聊天记录中，最新一次发送的消息id
		VersionID: versionId, // 2. 版本号（用于拉取会话框）
		SortKey:   versionId, // 3. sort_key的值等同于version_id
	}

	// save to db
	var res gen.ResultInfo
	res, err = qModel.Where(qModel.ID.Eq(contactId)).Limit(1).Updates(row)
	if err != nil {
		logging.Errorf(logHead+"Updates fail,err=%v,contact=%v", err, row)
		return
	}
	logging.Infof(logHead+"Updates success,contact=%v,RowsAffected=%v", row, res.RowsAffected)

	return
}

func (repo *ContactRepo) RangeList(logHead string, params *model.FetchContactRangeParams) (list []*model.Contact, err error) {
	logHead += "RangeList|"

	_, tbName := repo.TableName(params.OwnerId.Id())
	qModel := repo.Db.Contact.Table(tbName)
	pivotVersionId := params.PivotVersionId
	ownerId := params.OwnerId

	// 需要建立索引：owner_id、status、version_id
	switch params.FetchType {
	case model.FetchTypeBackward: // 拉取历史消息，范围为：（负无穷, pivotVersionId）
		if pivotVersionId == 0 {
			pivotVersionId = math.MaxInt64
		}
		list, err = qModel.Where(
			qModel.OwnerID.Eq(ownerId.Id()), qModel.OwnerType.Eq(ownerId.Type()),
			qModel.VersionID.Lt(pivotVersionId),
		).Limit(params.Limit).Order(qModel.VersionID.Desc()).Find()
	case model.FetchTypeForward: // 拉取最新消息，范围为：（pivotVersionId, 正无穷）
		list, err = qModel.Where(
			qModel.OwnerID.Eq(ownerId.Id()), qModel.OwnerType.Eq(ownerId.Type()),
			qModel.VersionID.Gt(pivotVersionId),
		).Limit(params.Limit).Order(qModel.VersionID).Find()
	}
	if err != nil {
		logging.Errorf(logHead+"err=%v", err)
		return
	}

	return
}

// UpdateLastDelMsg 更新contact的最后一条已删除的消息（清空聊天记录）
func (repo *ContactRepo) UpdateLastDelMsg(logHead string, lastDelMsgId model.BigIntType, versionId uint64, ownerId *gen_id.ComponentId, peerId *gen_id.ComponentId) (affectedRow int64, err error) {
	_, tbName := repo.TableName(ownerId.Id())
	qModel := repo.Db.Contact.Table(tbName)

	var res gen.ResultInfo
	res, err = qModel.Where(
		qModel.OwnerID.Eq(ownerId.Id()), qModel.OwnerType.Eq(ownerId.Type()),
		qModel.PeerID.Eq(peerId.Id()), qModel.PeerType.Eq(peerId.Type()),
	).Limit(1).Updates(&model.Contact{
		LastDelMsgID: lastDelMsgId,
		VersionID:    versionId,
	})
	if err != nil {
		logging.Errorf(logHead+"Update err=%v", err)
		return
	}
	affectedRow = res.RowsAffected

	return
}

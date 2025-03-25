package data

import (
	"context"
	"errors"
	"github.com/zhixunjie/im-fun/internal/logic/data/cache"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
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

// InfoWithCache 查询某个会话的信息
func (repo *ContactRepo) InfoWithCache(ownerId *gen_id.ComponentId, peerId *gen_id.ComponentId) (*model.Contact, error) {
	// todo 先从cache拿，拿不到再从DB拿

	return repo.Info("InfoWithCache|", ownerId, peerId)
}

// Info 查询某个会话的信息
func (repo *ContactRepo) Info(logHead string, ownerId *gen_id.ComponentId, peerId *gen_id.ComponentId) (row *model.Contact, err error) {
	dbName, tbName := model.ShardingTbNameContact(ownerId.Id())
	slave := repo.Slave(dbName).Contact.Table(tbName)

	row, err = slave.Where(
		slave.OwnerID.Eq(ownerId.Id()),
		slave.OwnerType.Eq(uint32(ownerId.Type())),
		slave.PeerID.Eq(peerId.Id()),
		slave.PeerType.Eq(uint32(peerId.Type())),
	).Take()
	if err != nil {
		return
	}
	return
}

// Edit 插入/更新记录
//func (repo *ContactRepo) Edit(logHead string, tx *query.Query, row *model.Contact) (err error) {
//	logHead += fmt.Sprintf("Edit,row=%v|", row)
//	dbName, tbName := model.ShardingTbNameContact(row.OwnerID)
//	master := repo.Master(dbName).Contact.Table(tbName)
//
//	// insert or update ?
//	if row.ID == 0 {
//		err = master.Create(row)
//		if err != nil {
//			logging.Errorf(logHead+"Create fail,err=%v", err)
//			return
//		}
//		logging.Infof(logHead + "Create success")
//	} else {
//		var res gen.ResultInfo
//		res, err = master.Where(master.ID.Eq(row.ID)).Limit(1).Updates(row)
//		if err != nil {
//			logging.Errorf(logHead+"Updates fail,err=%v", err)
//			return
//		}
//		logging.Infof(logHead+"Updates success,RowsAffected=%v", res.RowsAffected)
//	}
//
//	return
//}

// CreateNotExists 创建会话
func (repo *ContactRepo) CreateNotExists(logHead string, params *model.BuildContactParams) (contact *model.Contact, err error) {
	logHead += "CreateNotExists|"
	dbName, tbName := model.ShardingTbNameContact(params.Owner.Id())
	master := repo.Master(dbName).Contact.Table(tbName)

	// query contact
	contact, tmpErr := repo.Info(logHead, params.Owner, params.Peer)
	if tmpErr != nil && !errors.Is(tmpErr, gorm.ErrRecordNotFound) {
		err = tmpErr
		logging.Errorf(logHead+"Info error=%v", err)
		return
	}

	// insert if not exists
	if errors.Is(tmpErr, gorm.ErrRecordNotFound) {
		contact = &model.Contact{
			OwnerID:   params.Owner.Id(),
			OwnerType: uint32(params.Owner.Type()),
			PeerID:    params.Peer.Id(),
			PeerType:  uint32(params.Peer.Type()),
			Status:    model.ContactStatusNormal,
		}

		// save to db
		err = master.Create(contact)
		if err != nil {
			logging.Errorf(logHead+"Create fail,err=%v,contact=%v", err, contact)
			return
		}
		logging.Infof(logHead+"Create success,contact=%v", contact)
	}
	return
}

// RangeList 获取一定范围的会话列表
func (repo *ContactRepo) RangeList(logHead string, params *model.FetchContactRangeParams) (list []*model.Contact, err error) {
	logHead += "RangeList|"
	dbName, tbName := model.ShardingTbNameContact(params.Owner.Id())
	slave := repo.Slave(dbName).Contact.Table(tbName)

	pivotVersionId := params.PivotVersionId
	ownerId := params.Owner

	// 需要建立索引：owner_id、status、version_id
	switch params.FetchType {
	case model.FetchTypeBackward: // 拉取历史消息，范围为：（负无穷, pivotVersionId）
		if pivotVersionId == 0 {
			pivotVersionId = math.MaxInt64
		}
		list, err = slave.Where(
			slave.OwnerID.Eq(ownerId.Id()),
			slave.OwnerType.Eq(uint32(ownerId.Type())),
			slave.VersionID.Lt(pivotVersionId),
		).Limit(params.Limit).Order(slave.VersionID.Desc()).Find()
	case model.FetchTypeForward: // 拉取最新消息，范围为：（pivotVersionId, 正无穷）
		list, err = slave.Where(
			slave.OwnerID.Eq(ownerId.Id()),
			slave.OwnerType.Eq(uint32(ownerId.Type())),
			slave.VersionID.Gt(pivotVersionId),
		).Limit(params.Limit).Order(slave.VersionID).Find()
	default:
		err = errors.New("invalid fetch type")
		return
	}
	if err != nil {
		logging.Errorf(logHead+"err=%v", err)
		return
	}

	return
}

// UpdateLastMsgId 更新contact的最后一条消息（发消息）
func (repo *ContactRepo) UpdateLastMsgId(ctx context.Context, logHead string, contactId uint64, owner *gen_id.ComponentId, lastMsgId uint64, peerAck model.PeerAckStatus) (err error) {
	logHead += "UpdateLastMsgId|"
	mem := repo.RedisClient
	dbName, tbName := model.ShardingTbNameContact(owner.Id())
	master := repo.Master(dbName).Contact.Table(tbName)

	// note: 同一用户的会话timeline的版本变动，需要加锁
	lockKey := cache.TimelineContactLock.Format(k.M{"contact_id": contactId})
	redisSpinLock := distrib_lock.NewSpinLock(mem, lockKey, 5*time.Second, &distrib_lock.SpinOption{Interval: 20 * time.Millisecond, Times: 20})
	if err = redisSpinLock.AcquireWithTimes(); err != nil {
		logging.Errorf(logHead+"acquire fail,lockKey=%v,err=%v", lockKey, err)
		return
	}
	defer func() { _ = redisSpinLock.Release() }()
	logging.Infof(logHead+"acquire success,lockKey=%v", lockKey)

	// contact: get version_id
	versionId, err := gen_id.ContactVersionId(ctx, &gen_id.ContactVerParams{Mem: mem, Owner: owner})
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

	if peerAck == model.PeerAcked {
		row.PeerAck = uint32(peerAck) // 对方是否回应Owner
	}

	// save to db（要求：数据库的最后一条消息id小于当前消息id）
	res, err := master.Where(master.ID.Eq(contactId), master.LastMsgID.Lt(lastMsgId)).Limit(1).Updates(row)
	if err != nil {
		logging.Errorf(logHead+"Updates fail,err=%v,contact=%v", err, row)
		return
	}
	logging.Infof(logHead+"Updates success,contact=%v,RowsAffected=%v", row, res.RowsAffected)

	return
}

// UpdateLastDelMsg 更新contact的最后一条已删除的消息（清空聊天记录）
func (repo *ContactRepo) UpdateLastDelMsg(logHead string, lastDelMsgId model.BigIntType, versionId uint64, ownerId *gen_id.ComponentId, peerId *gen_id.ComponentId) (affectedRow int64, err error) {
	dbName, tbName := model.ShardingTbNameContact(ownerId.Id())
	master := repo.Master(dbName).Contact.Table(tbName)

	var res gen.ResultInfo
	res, err = master.Where(
		master.OwnerID.Eq(ownerId.Id()),
		master.OwnerType.Eq(uint32(ownerId.Type())),
		master.PeerID.Eq(peerId.Id()),
		master.PeerType.Eq(uint32(peerId.Type())),
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

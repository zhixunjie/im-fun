package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
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
func (repo *ContactRepo) InfoWithCache(ownerId, peerId *gmodel.ComponentId) (*model.ChatContact, error) {
	// todo 先从cache拿，拿不到再从DB拿

	return repo.Info(ownerId, peerId)
}

// Info 查询某个会话的信息
func (repo *ContactRepo) Info(ownerId, peerId *gmodel.ComponentId) (row *model.ChatContact, err error) {
	dbName, tbName := model.TbNameContact(ownerId.GetId())
	slave := repo.Slave(dbName).ChatContact.Table(tbName)

	row, err = slave.Where(
		slave.OwnerID.Eq(ownerId.GetId()),
		slave.OwnerType.Eq(uint32(ownerId.GetType())),
		slave.PeerID.Eq(peerId.GetId()),
		slave.PeerType.Eq(uint32(peerId.GetType())),
	).Take()
	if err != nil {
		return
	}
	return
}

// Edit 插入/更新记录
//func (repo *ContactRepo) Edit( tx *query.Query, row *model.ChatContact) (err error) {
//	logHead += fmt.Sprintf("Edit,row=%v|", row)
//	dbName, tbName := model.TbNameContact(row.OwnerID)
//	master := repo.Master(dbName).ChatContact.Table(tbName)
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
func (repo *ContactRepo) CreateNotExists(ctx context.Context, logHead string, params *model.BuildContactParams) (contact *model.ChatContact, err error) {
	logHead += "CreateNotExists|"
	dbName, tbName := model.TbNameContact(params.Owner.GetId())
	master := repo.Master(dbName).ChatContact.Table(tbName)

	// TODO: 使用 redis hash/string 进行优化（支持同时查两个contact）
	// query contact
	contact, tmpErr := repo.Info(params.Owner, params.Peer)
	if tmpErr != nil && !errors.Is(tmpErr, gorm.ErrRecordNotFound) {
		err = tmpErr
		logging.Errorf(logHead+"Info error=%v", err)
		return
	}

	// insert if not exists
	if errors.Is(tmpErr, gorm.ErrRecordNotFound) {
		contact = &model.ChatContact{
			OwnerID:   params.Owner.GetId(),
			OwnerType: uint32(params.Owner.GetType()),
			PeerID:    params.Peer.GetId(),
			PeerType:  uint32(params.Peer.GetType()),
			PeerAck:   uint32(gmodel.PeerNotAck),
			Status:    uint32(gmodel.ContactStatusNormal),
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
func (repo *ContactRepo) RangeList(params *model.FetchContactRangeParams) (list []*model.ChatContact, err error) {
	dbName, tbName := model.TbNameContact(params.Owner.GetId())
	slave := repo.Slave(dbName).ChatContact.Table(tbName)

	pivotVersionId := params.PivotVersionId
	ownerId := params.Owner

	// 需要建立索引：owner_id、status、version_id
	switch params.FetchType {
	case gmodel.FetchTypeBackward: // 拉取历史消息，范围为：（负无穷, pivotVersionId）
		if pivotVersionId == 0 {
			pivotVersionId = math.MaxInt64
		}
		list, err = slave.Where(
			slave.OwnerID.Eq(ownerId.GetId()),
			slave.OwnerType.Eq(uint32(ownerId.GetType())),
			slave.VersionID.Lt(pivotVersionId),
		).Limit(params.Limit).Order(slave.VersionID.Desc()).Find()
		if err != nil {
			err = fmt.Errorf("FetchTypeBackward err=%v", err)
			return
		}
	case gmodel.FetchTypeForward: // 拉取最新消息，范围为：（pivotVersionId, 正无穷）
		list, err = slave.Where(
			slave.OwnerID.Eq(ownerId.GetId()),
			slave.OwnerType.Eq(uint32(ownerId.GetType())),
			slave.VersionID.Gt(pivotVersionId),
		).Limit(params.Limit).Order(slave.VersionID).Find()
		if err != nil {
			err = fmt.Errorf("FetchTypeForward err=%v", err)
			return
		}
	default:
		err = errors.New("invalid fetch type")
		return
	}

	return
}

// UpdateLastMsgId 更新contact的最后一条消息（发消息）
func (repo *ContactRepo) UpdateLastMsgId(ctx context.Context, logHead string, contactId uint64, owner *gmodel.ComponentId, lastMsgId uint64, peerAck gmodel.PeerAckStatus) (err error) {
	logHead += "UpdateLastMsgId|"
	mem := repo.RedisClient
	dbName, tbName := model.TbNameContact(owner.GetId())
	master := repo.Master(dbName).ChatContact.Table(tbName)

	// note: 同一用户的会话timeline的版本变动，需要加锁
	lockKey := TimelineContactLock.Format(k.M{"contact_id": contactId})
	redisSpinLock := distrib_lock.NewSpinLock(mem, lockKey, 5*time.Second, &distrib_lock.SpinOption{Interval: 50 * time.Millisecond, Times: 40})
	if err = redisSpinLock.AcquireWithTimes(); err != nil {
		logging.Errorf(logHead+"acquire fail,lockKey=%v,err=%v", lockKey, err)
		return
	}
	defer func() { _ = redisSpinLock.Release() }()
	logging.Infof(logHead+"acquire success,lockKey=%v", lockKey)

	// generate contact's version_id
	versionId, err := gen_id.NewContactVersionId(ctx, &gen_id.ContactVerParams{Mem: mem, Owner: owner})
	if err != nil {
		logging.Errorf(logHead+"gen VersionID error=%v", err)
		return
	}

	// 只更新一部分的字段
	row := &model.ChatContact{
		LastMsgID: lastMsgId, // 1. 双方聊天记录中，最新一次发送的消息id
		VersionID: versionId, // 2. 版本号（用于拉取会话框）
		SortKey:   versionId, // 3. sort_key的值等同于version_id
	}

	if peerAck == gmodel.PeerAcked {
		row.PeerAck = uint32(peerAck) // 对方是否回应Owner
	}

	// 要求：数据库的 last_msg_id 小于当前的 last_msg_id
	res, err := master.Where(master.ID.Eq(contactId), master.LastMsgID.Lt(lastMsgId)).Limit(1).Updates(row)
	if err != nil {
		logging.Errorf(logHead+"Updates fail,err=%v,contact=%v", err, row)
		return
	}
	logging.Infof(logHead+"Updates success,contact=%v,RowsAffected=%v", row, res.RowsAffected)

	return
}

// UpdateLastDelMsg 更新contact的最后一条已删除的消息（清空聊天记录）
func (repo *ContactRepo) UpdateLastDelMsg(lastDelMsgId model.BigIntType, versionId uint64, ownerId, peerId *gmodel.ComponentId) (err error) {
	dbName, tbName := model.TbNameContact(ownerId.GetId())
	master := repo.Master(dbName).ChatContact.Table(tbName)

	var res gen.ResultInfo
	res, err = master.Where(
		master.OwnerID.Eq(ownerId.GetId()),
		master.OwnerType.Eq(uint32(ownerId.GetType())),
		master.PeerID.Eq(peerId.GetId()),
		master.PeerType.Eq(uint32(peerId.GetType())),
	).Limit(1).Updates(&model.ChatContact{
		LastDelMsgID: lastDelMsgId,
		VersionID:    versionId,
	})
	if err != nil {
		return
	}
	if res.RowsAffected == 0 {
		err = errors.New("affectedRow not allow")
		return
	}

	return
}

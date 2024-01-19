package data

import (
	"context"
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/query"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"gorm.io/gen"
	"gorm.io/gorm"
	"math"
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

	return repo.Info("InfoWithCache|", ownerId, peerId)
}

// Info 查询某个会话的信息
func (repo *ContactRepo) Info(logHead string, ownerId uint64, peerId uint64) (row *model.Contact, err error) {
	_, tbName := repo.TableName(ownerId)
	qModel := repo.Db.Contact.Table(tbName)

	row, err = qModel.Where(qModel.OwnerID.Eq(ownerId), qModel.PeerID.Eq(peerId)).Take()
	if err != nil {
		logging.Errorf(logHead+"Take err=%v", err)
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

func (repo *ContactRepo) Build(ctx context.Context, logHead string, params *model.BuildContactParams) (contact *model.Contact, err error) {
	logHead += "Build|"
	mem := repo.RedisClient

	// get version_id
	versionId, err := gen_id.VersionId(ctx, &gen_id.GenVersionParams{
		Mem:            mem,
		GenVersionType: gen_id.GenVersionTypeContact,
		OwnerId:        params.OwnerId,
	})
	if err != nil {
		logging.Errorf(logHead+"gen VersionId error=%v", err)
		return
	}

	// query contact
	contact, err = repo.Info(logHead, params.OwnerId, params.PeerId)
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.Errorf(logHead+"Info error=%v", err)
		return
	}

	// 记录不存在：需要创建contact
	if err == gorm.ErrRecordNotFound {
		err = nil
		contact = &model.Contact{
			OwnerID: params.OwnerId,
			PeerID:  params.PeerId,
			//ContactPeerType:  params.ContactPeerType, // 注意：暂时不适用请求参数过来的PeerType（适合于：logic -> base的场景）
			PeerType:  int32(model.PeerTypeNormalUser),
			PeerAck:   params.InitPeerAck,
			LastMsgID: params.LastMsgId,
			VersionID: versionId,
			SortKey:   versionId,
			Status:    model.ContactStatusNormal,
		}
	} else {
		contact.LastMsgID = params.LastMsgId // 双方聊天记录中，最新一次发送的消息id
		contact.VersionID = versionId        // 版本号（用于拉取会话框）
		contact.SortKey = versionId          // sort_key的值等同于version_id
	}
	return
}

func (repo *ContactRepo) RangeList(params *model.FetchContactRangeParams) (list []*model.Contact, err error) {
	logHead := "RangeList|"

	_, tbName := repo.TableName(params.OwnerId)
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
			qModel.OwnerID.Eq(ownerId),
			qModel.Status.Eq(uint32(model.ContactStatusNormal)),
			qModel.VersionID.Lt(pivotVersionId),
		).Limit(params.Limit).Order(qModel.VersionID.Desc()).Find()
	case model.FetchTypeForward: // 拉取最新消息，范围为：（pivotVersionId, 正无穷）
		list, err = qModel.Where(
			qModel.OwnerID.Eq(ownerId),
			qModel.Status.Eq(uint32(model.ContactStatusNormal)),
			qModel.VersionID.Gt(pivotVersionId),
		).Limit(params.Limit).Order(qModel.VersionID).Find()
	}
	if err != nil {
		logging.Errorf(logHead+"err=%v", err)
		return
	}

	return
}

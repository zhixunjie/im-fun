package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"gorm.io/gen"
	"math"
	"strings"
)

type GroupMessageRepo struct {
	*Data
}

func NewGroupMessageRepo(data *Data) *GroupMessageRepo {
	return &GroupMessageRepo{
		Data: data,
	}
}

func (repo *GroupMessageRepo) Create(logHead string, row *model.ChatGroupMessage) (err error) {
	logHead += "Create|"
	dbName, tbName := model.TbNameGroupMessage(row.MsgID)
	master := repo.Master(dbName).ChatGroupMessage.Table(tbName)

	err = master.Create(row)
	if err != nil {
		logging.Errorf(logHead+"Create fail err=%v", err)
		return
	}
	logging.Infof(logHead+"Create success(%v,%v),rowId=%v", dbName, tbName, row.ID)

	return
}

// InfoWithCache 查询某条消息的详情
func (repo *GroupMessageRepo) InfoWithCache(msgId uint64) (*model.ChatGroupMessage, error) {
	// todo 先从cache拿，拿不到再从DB拿

	return repo.Info(msgId)
}

// Info 查询某条消息的详情
func (repo *GroupMessageRepo) Info(msgId uint64) (row *model.ChatGroupMessage, err error) {
	dbName, tbName := model.TbNameGroupMessage(msgId)
	slave := repo.Slave(dbName).ChatGroupMessage.Table(tbName)

	row, err = slave.Where(slave.MsgID.Eq(msgId)).Take()
	if err != nil {
		return
	}
	return
}

// RangeList 获取一定范围的消息列表
func (repo *GroupMessageRepo) RangeList(params *model.FetchMsgRangeParams) (list []*model.ChatGroupMessage, err error) {
	dbName, tbName := model.TbNameGroupMessageByCId(params.Owner, params.Peer)
	slave := repo.Slave(dbName).ChatGroupMessage.Table(tbName)

	// get id
	delVersionId := params.LastDelMsgVersionId
	pivotVersionId := params.PivotVersionId

	// 需要建立索引：session_id、status、version_id
	switch params.FetchType {
	case gmodel.FetchTypeBackward: // 📚拉取历史消息，范围为：（delVersionId, pivotVersionId）
		if pivotVersionId == 0 {
			pivotVersionId = math.MaxInt64
		}
		list, err = slave.Where(
			slave.SessionID.Eq(string(params.SessionId)),
			slave.VersionID.Gt(delVersionId),
			slave.VersionID.Lt(pivotVersionId),
		).Limit(params.Limit).Order(slave.VersionID.Desc()).Find() // 按照version_id倒序排序
		if err != nil {
			err = fmt.Errorf("FetchTypeBackward err=%v", err)
			return
		}
	case gmodel.FetchTypeForward: // 📚拉取最新消息，范围为：（pivotVersionId, 正无穷）
		// 避免：拉取最新消息时拉到已删除消息
		if pivotVersionId < delVersionId {
			pivotVersionId = delVersionId
		}
		list, err = slave.Where(
			slave.SessionID.Eq(string(params.SessionId)),
			slave.VersionID.Gt(pivotVersionId),
		).Limit(params.Limit).Order(slave.VersionID).Find() // 按照version_id正序排序
		if err != nil {
			err = fmt.Errorf("FetchTypeForward err=%v", err)
			return
		}
	default:
		err = errors.New("invalid fetchType")
		return
	}

	return
}

// UpdateMsgVerAndStatus 修改某条消息的状态
func (repo *GroupMessageRepo) UpdateMsgVerAndStatus(logHead string, msgId, versionId model.BigIntType, status gmodel.MsgStatus) (err error) {
	logHead += fmt.Sprintf("UpdateMsgVerAndStatus,msgId=%v,versionId=%v,status=%v|", msgId, versionId, status)
	dbName, tbName := model.TbNameGroupMessage(msgId)
	master := repo.Master(dbName).ChatGroupMessage.Table(tbName)

	// status
	srcStatus := uint32(gmodel.MsgStatusNormal)
	dstStatus := uint32(status)

	// operation
	var res gen.ResultInfo
	res, err = master.
		Where(master.MsgID.Eq(msgId), master.Status.Eq(srcStatus)).Limit(1).
		Updates(&model.ChatGroupMessage{
			VersionID: versionId,
			Status:    dstStatus,
		})
	if err != nil {
		logging.Errorf(logHead+"UpdateMsgVerAndStatus error=%v", err)
		return
	}
	if res.RowsAffected == 0 {
		err = errors.New("affectedRow not allow")
		logging.Errorf(logHead+"UpdateMsgVerAndStatus error=%v", err)
		return
	}

	return
}

func (repo *GroupMessageRepo) BatchGetByMsgIds(ctx context.Context, msgIds []uint64) (retMap map[uint64]*model.ChatGroupMessage, err error) {
	tbNames := map[string][]uint64{}
	for _, msgId := range msgIds {
		dbName, tbName := model.TbNameGroupMessage(msgId)
		key := dbName + "_" + tbName
		tbNames[key] = append(tbNames[key], msgId)
	}

	list := make([]*model.ChatGroupMessage, 0, len(msgIds))
	tmpList := make([]*model.ChatGroupMessage, 0, len(msgIds))
	for key, ids := range tbNames {
		res := strings.Split(key, "_")
		if len(res) != 2 {
			continue
		}
		dbName, tbName := res[0], res[1]
		slave := repo.Slave(dbName).ChatGroupMessage.Table(tbName)
		// find it
		tmpList, err = slave.Where(slave.MsgID.In(ids...)).Find()
		if err != nil {
			err = fmt.Errorf("BatchGetByMsgIds Find in err=%v", err)
			return
		}
		list = append(list, tmpList...)
	}
	return
}

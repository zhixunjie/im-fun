package data

import (
	"errors"
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/pkg/logging"
	"gorm.io/gen"
	"math"
)

type MessageRepo struct {
	*Data
}

func NewMessageRepo(data *Data) *MessageRepo {
	return &MessageRepo{
		Data: data,
	}
}

func (repo *MessageRepo) Create(logHead string, row *model.Message) (err error) {
	logHead += "Create|"
	dbName, tbName := model.ShardingTbNameMessage(row.MsgID)
	master := repo.Master(dbName).Message.Table(tbName)

	err = master.Create(row)
	if err != nil {
		logging.Errorf(logHead+"Create fail err=%v", err)
		return
	}
	logging.Infof(logHead+"Create success,rowId=%v", row.ID)

	return
}

// InfoWithCache 查询某条消息的详情
func (repo *MessageRepo) InfoWithCache(msgId uint64) (*model.Message, error) {
	// todo 先从cache拿，拿不到再从DB拿

	return repo.Info(msgId)
}

// Info 查询某条消息的详情
func (repo *MessageRepo) Info(msgId uint64) (row *model.Message, err error) {
	dbName, tbName := model.ShardingTbNameMessage(msgId)
	slave := repo.Slave(dbName).Message.Table(tbName)

	row, err = slave.Where(slave.MsgID.Eq(msgId)).Take()
	if err != nil {
		return
	}
	return
}

// RangeList 获取一定范围的消息列表
func (repo *MessageRepo) RangeList(params *model.FetchMsgRangeParams) (list []*model.Message, err error) {
	dbName, tbName := model.ShardingTbNameMessageByComponentId(params.OwnerId, params.PeerId)
	slave := repo.Slave(dbName).Message.Table(tbName)

	// get id
	delVersionId := params.LastDelMsgVersionId
	pivotVersionId := params.PivotVersionId

	// 需要建立索引：session_id、status、version_id
	switch params.FetchType {
	case model.FetchTypeBackward: // 拉取历史消息，范围为：（delVersionId, pivotVersionId）
		if pivotVersionId == 0 {
			pivotVersionId = math.MaxInt64
		}
		list, err = slave.Where(
			slave.SessionID.Eq(params.SessionId),
			slave.VersionID.Gt(delVersionId),
			slave.VersionID.Lt(pivotVersionId),
		).Limit(params.Limit).Order(slave.VersionID.Desc()).Find() // 按照version_id倒序排序
	case model.FetchTypeForward: // 拉取最新消息，范围为：（pivotVersionId, 正无穷）
		// 避免：拉取最新消息时拉到已删除消息
		if pivotVersionId < delVersionId {
			pivotVersionId = delVersionId
		}
		list, err = slave.Where(
			slave.SessionID.Eq(params.SessionId),
			slave.VersionID.Gt(pivotVersionId),
		).Limit(params.Limit).Order(slave.VersionID).Find() // 按照version_id正序排序
	default:
		err = errors.New("invalid fetchType")
		return
	}

	return
}

// UpdateMsgVerAndStatus 修改某条消息的状态
func (repo *MessageRepo) UpdateMsgVerAndStatus(logHead string, msgId, versionId model.BigIntType, status model.MsgStatus) (err error) {
	logHead += fmt.Sprintf("UpdateMsgVerAndStatus,msgId=%v,versionId=%v,status=%v|", msgId, versionId, status)
	dbName, tbName := model.ShardingTbNameMessage(msgId)
	master := repo.Master(dbName).Message.Table(tbName)

	// status
	srcStatus := uint32(model.MsgStatusNormal)
	dstStatus := uint32(status)

	// operation
	var res gen.ResultInfo
	res, err = master.
		Where(master.MsgID.Eq(msgId), master.Status.Eq(srcStatus)).Limit(1).
		Updates(&model.Message{
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

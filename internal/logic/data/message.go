package data

import (
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/query"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
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

// TableName
// 因为msgId和largerId的后4位是相同的，所以这里传入msgId或者largerId都可以
func (repo *MessageRepo) TableName(id uint64) (dbName string, tbName string) {
	// 临时写死
	if true {
		return "", "message"
	}
	// 分表规则：
	// - 数据库前缀：message_xxx，规则：id 倒数第三位数字就是分库值
	// - 数据表前缀：message_xxx，规则：id 的最后两位就是分表值
	dbName = fmt.Sprintf("messsage_%v", id%1000/100)
	tbName = fmt.Sprintf("message_%v", id%model.TotalTableMessage)

	return dbName, tbName
}

func (repo *MessageRepo) Create(logHead string, tx *query.Query, row *model.Message) (err error) {
	logHead += "Create|"
	_, tbName := repo.TableName(row.MsgID)
	qModel := tx.Message

	err = qModel.Table(tbName).Create(row)
	if err != nil {
		logging.Errorf(logHead+"Create fail err=%v", err)
		return
	}
	logging.Infof(logHead+"Create success,row=%v", row)

	return
}

// InfoWithCache 查询某条消息的详情
func (repo *MessageRepo) InfoWithCache(msgId uint64) (*model.Message, error) {
	// todo 先从cache拿，拿不到再从DB拿

	return repo.Info(msgId)
}

// Info 查询某条消息的详情
func (repo *MessageRepo) Info(msgId uint64) (row *model.Message, err error) {
	_, tbName := repo.TableName(msgId)
	qModel := repo.Db.Message.Table(tbName)

	row, err = qModel.Where(qModel.MsgID.Eq(msgId)).Take()
	if err != nil {
		return
	}
	return
}

func (repo *MessageRepo) RangeList(params *model.FetchMsgRangeParams) (list []*model.Message, err error) {
	_, tbName := repo.TableName(params.LargerId.Id())
	qModel := repo.Db.Message.Table(tbName)

	// get id
	sessionId := gen_id.GenSessionId(params.SmallerId, params.LargerId)
	delVersionId := params.LastDelMsgVersionId
	pivotVersionId := params.PivotVersionId

	// 需要建立索引：session_id、status、version_id
	switch params.FetchType {
	case model.FetchTypeBackward: // 拉取历史消息，范围为：（delVersionId, pivotVersionId）
		if pivotVersionId == 0 {
			pivotVersionId = math.MaxInt64
		}
		list, err = qModel.Where(
			qModel.SessionID.Eq(sessionId),
			qModel.VersionID.Gt(delVersionId),
			qModel.VersionID.Lt(pivotVersionId),
		).Limit(params.Limit).Order(qModel.VersionID.Desc()).Find() // 按照version_id倒序排序
	case model.FetchTypeForward: // 拉取最新消息，范围为：（pivotVersionId, 正无穷）
		// 避免：拉取最新消息时拉到已删除消息
		if pivotVersionId < delVersionId {
			pivotVersionId = delVersionId
		}
		list, err = qModel.Where(
			qModel.SessionID.Eq(sessionId),
			qModel.VersionID.Gt(pivotVersionId),
		).Limit(params.Limit).Order(qModel.VersionID).Find() // 按照version_id正序排序
	}

	return
}

func (repo *MessageRepo) UpdateMsgStatus(msgId model.BigIntType, status model.MsgStatus, versionId uint64) (affectedRow int64, err error) {
	_, tbName := repo.TableName(msgId)
	qModel := repo.Db.Message.Table(tbName)

	// status
	srcStatus := uint32(model.MsgStatusNormal)
	dstStatus := uint32(status)

	// operation
	var res gen.ResultInfo
	res, err = qModel.
		Where(qModel.Status.Eq(srcStatus)).Limit(1).
		Updates(&model.Message{
			VersionID: versionId,
			Status:    dstStatus,
		})
	if err != nil {
		return
	}
	affectedRow = res.RowsAffected

	return
}

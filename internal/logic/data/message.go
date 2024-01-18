package data

import (
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/query"
	"github.com/zhixunjie/im-fun/pkg/gen_id"
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
	//// 临时写死
	//if true {
	//	return "", "message_0"
	//}
	// 分表规则：
	// - 数据库前缀：message_xxx，规则：id 倒数第三位数字就是分库值
	// - 数据表前缀：message_xxx，规则：id 的最后两位就是分表值
	dbName = fmt.Sprintf("messsage_%v", id%1000/100)
	tbName = fmt.Sprintf("message_%v", id%model.TotalTableMessage)

	return dbName, tbName
}

func (repo *MessageRepo) Create(tx *query.Query, row *model.Message) (err error) {
	_, tbName := repo.TableName(row.MsgID)
	qModel := tx.Message

	err = qModel.Table(tbName).Create(row)
	if err != nil {
		return
	}
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

func (repo *MessageRepo) RangeList(params *model.QueryMsgParams) (list []*model.Message, err error) {
	_, tbName := repo.TableName(params.LargerId)
	qModel := repo.Db.Message.Table(tbName)
	sessionId := gen_id.SessionId(params.SmallerId, params.LargerId)
	delVersionId := params.DelVersionId
	pivotVersionId := params.PivotVersionId

	// 需要建立索引：session_id、status、version_id
	switch params.FetchType {
	case model.FetchTypeBackward: // 拉取历史消息，范围为：（delVersionId, pivotVersionId）
		if pivotVersionId == 0 {
			pivotVersionId = math.MaxInt64
		}
		list, err = qModel.Where(
			qModel.SessionID.Eq(sessionId),
			qModel.Status.Eq(model.MsgStatusNormal),
			qModel.VersionID.Gt(delVersionId),
			qModel.VersionID.Lt(pivotVersionId),
		).Limit(params.Limit).Order(qModel.VersionID.Desc()).Find()
	case model.FetchTypeForward: // 拉取最新消息，范围为：（pivotVersionId, 正无穷）
		list, err = qModel.Where(
			qModel.SessionID.Eq(sessionId),
			qModel.Status.Eq(model.MsgStatusNormal),
			qModel.VersionID.Gt(pivotVersionId),
		).Limit(params.Limit).Order(qModel.VersionID).Find()
	}

	return
}

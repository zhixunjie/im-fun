package data

import (
	"errors"
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
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

func (repo *MessageRepo) TableNameByContactId(id1, id2 *gen_id.ComponentId) (dbName string, tbName string) {
	switch {
	case id1.IsGroup(): // 群聊
		dbName, tbName = repo.TableName(id1.Id())
	case id2.IsGroup(): // 群聊
		dbName, tbName = repo.TableName(id2.Id())
	default: // 单聊
		_, largerId := gen_id.Sort(id1, id2)
		dbName, tbName = repo.TableName(largerId.Id())
	}

	return
}

// TableName
// 因为msgId和largerId的后4位是相同的，所以这里传入msgId或者largerId都可以
func (repo *MessageRepo) TableName(id uint64) (dbName string, tbName string) {
	// 临时写死
	if true {
		return "", "message"
	}
	// 分表规则：
	// - 数据库前缀：_xxx，       规则：id 倒数第四位数字就是分库值
	// - 数据表前缀：message_xxx，规则：id 的最后4位哈希分表
	// 🔥其实后四位都可以用来取余得到分表数，所有分表数是不止2位的
	last4bit := id % 1000 // 提取最后4位出来
	dbName = fmt.Sprintf("messsage_%v", last4bit/100)
	tbName = fmt.Sprintf("message_%v", last4bit%model.TotalTableMessage)

	return dbName, tbName
}

func (repo *MessageRepo) Create(logHead string, row *model.Message) (err error) {
	logHead += "Create|"
	_, tbName := repo.TableName(row.MsgID)
	qModel := repo.Db.Message.Table(tbName)

	err = qModel.Create(row)
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
	_, tbName := repo.TableName(msgId)
	qModel := repo.Db.Message.Table(tbName)

	row, err = qModel.Where(qModel.MsgID.Eq(msgId)).Take()
	if err != nil {
		return
	}
	return
}

// RangeList 获取一定范围的消息列表
func (repo *MessageRepo) RangeList(params *model.FetchMsgRangeParams) (list []*model.Message, err error) {
	_, tbName := repo.TableNameByContactId(params.OwnerId, params.PeerId)
	qModel := repo.Db.Message.Table(tbName)

	// get id
	delVersionId := params.LastDelMsgVersionId
	pivotVersionId := params.PivotVersionId

	// 需要建立索引：session_id、status、version_id
	switch params.FetchType {
	case model.FetchTypeBackward: // 拉取历史消息，范围为：（delVersionId, pivotVersionId）
		if pivotVersionId == 0 {
			pivotVersionId = math.MaxInt64
		}
		list, err = qModel.Where(
			qModel.SessionID.Eq(params.SessionId),
			qModel.VersionID.Gt(delVersionId),
			qModel.VersionID.Lt(pivotVersionId),
		).Limit(params.Limit).Order(qModel.VersionID.Desc()).Find() // 按照version_id倒序排序
	case model.FetchTypeForward: // 拉取最新消息，范围为：（pivotVersionId, 正无穷）
		// 避免：拉取最新消息时拉到已删除消息
		if pivotVersionId < delVersionId {
			pivotVersionId = delVersionId
		}
		list, err = qModel.Where(
			qModel.SessionID.Eq(params.SessionId),
			qModel.VersionID.Gt(pivotVersionId),
		).Limit(params.Limit).Order(qModel.VersionID).Find() // 按照version_id正序排序
	}

	return
}

// UpdateMsgVerAndStatus 修改某条消息的状态
func (repo *MessageRepo) UpdateMsgVerAndStatus(logHead string, msgId, versionId model.BigIntType, status model.MsgStatus) (err error) {
	logHead += fmt.Sprintf("UpdateMsgVerAndStatus,msgId=%v,versionId=%v,status=%v|", msgId, versionId, status)
	_, tbName := repo.TableName(msgId)
	qModel := repo.Db.Message.Table(tbName)

	// status
	srcStatus := uint32(model.MsgStatusNormal)
	dstStatus := uint32(status)

	// operation
	var res gen.ResultInfo
	res, err = qModel.
		Where(qModel.MsgID.Eq(msgId), qModel.Status.Eq(srcStatus)).Limit(1).
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

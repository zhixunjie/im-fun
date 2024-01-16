package dao

import (
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/model"
)

// GetMessageDbAndTable
// 因为msgId和largerId的后4位是相同的，所以这里传入msgId或者largerId都可以
func (d *Dao) GetMessageDbAndTable(id uint64) (dbName string, tbName string) {
	// 临时写死
	if true {
		return "", "message_0"
	}
	// 分表规则：
	// - 数据库前缀：message_xxx，规则：id 倒数第三位数字就是分库值
	// - 数据表前缀：message_xxx，规则：id 的最后两位就是分表值
	dbName = fmt.Sprintf("messsage_%v", id%1000/100)
	tbName = fmt.Sprintf("message_%v", id%model.TotalTableMessage)

	return dbName, tbName
}

// QueryMsgLogic 查询某条消息的详情
func (d *Dao) QueryMsgLogic(msgId uint64) (model.Message, error) {
	// todo 先从cache拿，拿不到再从DB拿

	return d.QueryMsgByMsgId(msgId)
}

// QueryMsgByMsgId 查询某条消息的详情，From：DB
func (d *Dao) QueryMsgByMsgId(msgId uint64) (model.Message, error) {
	_, tbName := d.GetMessageDbAndTable(msgId)
	var row model.Message
	err := d.MySQLClient.Table(tbName).Where("msg_id=?", msgId).Find(&row).Error
	if err != nil {
		return row, err
	}
	return row, nil
}

func (d *Dao) AddMsg(row *model.Message) error {
	_, tbName := d.GetMessageDbAndTable(row.MsgId)
	err := d.MySQLClient.Table(tbName).Create(&row).Error
	if err != nil {
		return err
	}
	return nil
}

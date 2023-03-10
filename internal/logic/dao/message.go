package dao

import (
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/model"
)

// GetMessageDbAndTable 这里传入msgId或者largerId都是可以的
func GetMessageDbAndTable(id uint64) (string, string) {
	// 临时写死
	if true {
		return "", "message_0"
	}
	// 分表规则
	// message表和contact表都放在message库
	dbName := fmt.Sprintf("messsage_%v", id%model.TotalDb)
	tbName := fmt.Sprintf("message_%v", id%model.TotalTableMessage)

	return dbName, tbName
}

// QueryMsgLogic 查询某条消息的详情
func QueryMsgLogic(msgId uint64) (model.Message, error) {
	// todo 先从cache拿，拿不到再从DB拿

	return QueryMsgByMsgId(msgId)
}

// QueryMsgByMsgId 查询某条消息的详情，From：DB
func QueryMsgByMsgId(msgId uint64) (model.Message, error) {
	db := MySQLClient
	_, tbName := GetMessageDbAndTable(msgId)
	var row model.Message
	err := db.Table(tbName).Where("msg_id=?", msgId).Find(&row).Error
	if err != nil {
		return row, err
	}
	return row, nil
}

func AddMsg(row *model.Message) error {
	db := MySQLClient
	_, tbName := GetMessageDbAndTable(row.MsgId)
	err := db.Table(tbName).Create(&row).Error
	if err != nil {
		return err
	}
	return nil
}

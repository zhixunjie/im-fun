package dao

import (
	"fmt"
	model2 "github.com/zhixunjie/im-fun/internal/logic/model"
	"gorm.io/gorm"
)

func GetContactDbAndTable(ownerId uint64) (string, string) {
	// 临时写死
	if true {
		return "", "contact_0"
	}
	// 分表规则
	// message表和contact表都放在message库
	dbName := fmt.Sprintf("messsage_%v", ownerId%model2.TotalDb)
	tbName := fmt.Sprintf("contact_%v", ownerId%model2.TotalTableMessage)

	return dbName, tbName
}

// QueryContactLogic 查询某个会话的信息
func QueryContactLogic(ownerId uint64, peerId uint64) (model2.Contact, error) {
	// todo 先从cache拿，拿不到再从DB拿

	return QueryContactById(ownerId, peerId)
}

// QueryContactById 查询某个会话的信息，From：DB
func QueryContactById(ownerId uint64, peerId uint64) (model2.Contact, error) {
	db := MySQLClient
	_, tbName := GetContactDbAndTable(ownerId)
	var row model2.Contact
	err := db.Table(tbName).Where("owner_id = ? AND peer_id = ?", ownerId, peerId).Find(&row).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return row, err
	}
	return row, nil
}

// AddOrUpdateContact 插入 or 更新记录
func AddOrUpdateContact(row *model2.Contact) error {
	db := MySQLClient
	_, tbName := GetContactDbAndTable(row.OwnerId)

	if row.Id == 0 {
		err := db.Table(tbName).Create(&row).Error
		if err != nil {
			return err
		}
	} else {
		err := db.Table(tbName).Updates(&row).Error
		if err != nil {
			return err
		}
	}

	return nil
}

package dao

import (
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/model"
	"gorm.io/gorm"
)

func (d *Dao) GetContactDbAndTable(ownerId uint64) (string, string) {
	// 临时写死
	if true {
		return "", "contact_0"
	}
	// 分表规则
	// message表和contact表都放在message库
	dbName := fmt.Sprintf("messsage_%v", ownerId%model.TotalDb)
	tbName := fmt.Sprintf("contact_%v", ownerId%model.TotalTableMessage)

	return dbName, tbName
}

// QueryContactLogic 查询某个会话的信息
func (d *Dao) QueryContactLogic(ownerId uint64, peerId uint64) (model.Contact, error) {
	// todo 先从cache拿，拿不到再从DB拿

	return d.QueryContactById(ownerId, peerId)
}

// QueryContactById 查询某个会话的信息，From：DB
func (d *Dao) QueryContactById(ownerId uint64, peerId uint64) (model.Contact, error) {
	_, tbName := d.GetContactDbAndTable(ownerId)
	var row model.Contact
	err := d.MySQLClient.Table(tbName).Where("owner_id = ? AND peer_id = ?", ownerId, peerId).Find(&row).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return row, err
	}
	return row, nil
}

// AddOrUpdateContact 插入 or 更新记录
func (d *Dao) AddOrUpdateContact(row *model.Contact) error {
	_, tbName := d.GetContactDbAndTable(row.OwnerId)

	if row.Id == 0 {
		err := d.MySQLClient.Table(tbName).Create(&row).Error
		if err != nil {
			return err
		}
	} else {
		err := d.MySQLClient.Table(tbName).Updates(&row).Error
		if err != nil {
			return err
		}
	}

	return nil
}

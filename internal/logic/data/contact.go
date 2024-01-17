package data

import (
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"gorm.io/gorm"
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
	// 临时写死
	if true {
		return "", "contact_0"
	}
	// 分表规则：
	// - 数据库前缀：message_xxx，规则：owner_id 倒数第三位数字就是分库值
	// - 数据表前缀：contact_xxx，规则：owner_id 的最后两位就是分表值
	dbName = fmt.Sprintf("messsage_%v", ownerId%1000/100)
	tbName = fmt.Sprintf("contact_%v", ownerId%model.TotalTableContact)

	return dbName, tbName
}

// QueryContactLogic 查询某个会话的信息
func (repo *ContactRepo) QueryContactLogic(ownerId uint64, peerId uint64) (model.Contact, error) {
	// todo 先从cache拿，拿不到再从DB拿

	return repo.QueryContactById(ownerId, peerId)
}

// QueryContactById 查询某个会话的信息，From：DB
func (repo *ContactRepo) QueryContactById(ownerId uint64, peerId uint64) (row model.Contact, err error) {
	_, tbName := repo.TableName(ownerId)
	err = repo.MySQLClient.Table(tbName).Where("owner_id = ? AND peer_id = ?", ownerId, peerId).Find(&row).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return row, err
	}
	return row, nil
}

// AddOrUpdateContact 插入 or 更新记录
func (repo *ContactRepo) AddOrUpdateContact(row *model.Contact) error {
	_, tbName := repo.TableName(row.OwnerID)

	if row.ID == 0 {
		err := repo.MySQLClient.Table(tbName).Create(&row).Error
		if err != nil {
			return err
		}
	} else {
		err := repo.MySQLClient.Table(tbName).Updates(&row).Error
		if err != nil {
			return err
		}
	}

	return nil
}

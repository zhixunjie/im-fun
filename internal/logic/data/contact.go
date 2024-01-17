package data

import (
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/query"
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
func (repo *ContactRepo) QueryContactLogic(ownerId uint64, peerId uint64) (*model.Contact, error) {
	// todo 先从cache拿，拿不到再从DB拿

	return repo.QueryContactById(ownerId, peerId)
}

// QueryContactById 查询某个会话的信息，From：DB
func (repo *ContactRepo) QueryContactById(ownerId uint64, peerId uint64) (row *model.Contact, err error) {
	_, tbName := repo.TableName(ownerId)
	qModel := repo.Db.Contact

	row, err = qModel.Table(tbName).Where(
		qModel.OwnerID.Eq(ownerId),
		qModel.PeerID.Eq(peerId),
	).Take()
	if err != nil {
		return
	}
	return
}

// AddOrUpdateContact 插入 or 更新记录
func (repo *ContactRepo) AddOrUpdateContact(tx *query.Query, row *model.Contact) (err error) {
	_, tbName := repo.TableName(row.OwnerID)
	qModel := tx.Contact

	if row.ID == 0 {
		err = qModel.Table(tbName).Create(row)
		if err != nil {
			return
		}
	} else {
		_, err = qModel.Table(tbName).Updates(row)
		if err != nil {
			return
		}
	}

	return
}

package data

import (
	"fmt"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

type UserRepo struct {
	*Data
}

func NewUserRepo(data *Data) *UserRepo {
	return &UserRepo{
		Data: data,
	}
}

func (repo *UserRepo) GetUserByAccount(accountType uint32, accountId string) (item *model.User, err error) {
	item = new(model.User)
	qModel := repo.MySQLDb.User

	item, err = qModel.Where(
		qModel.AccountType.Eq(accountType),
		qModel.AccountID.Eq(accountId),
	).Take()
	if err != nil {
		err = fmt.Errorf("CreateUser fail err=%w", err)
		return
	}
	return
}

func (repo *UserRepo) CreateUser(item *model.User) (err error) {
	qModel := repo.MySQLDb.User
	err = qModel.Create(item)
	if err != nil {
		err = fmt.Errorf("CreateUser fail err=%v", err)
		return
	}
	return
}

func (repo *UserRepo) DeleteUser(id uint64) (err error) {
	qModel := repo.MySQLDb.User
	res, err := qModel.Where(qModel.ID.Eq(id)).Delete()
	if err != nil {
		logging.Errorf("DeleteUser err=%v", err)
		return
	}
	logging.Infof("DeleteUser res=%v", res)
	return
}

package biz

import "github.com/zhixunjie/im-fun/internal/logic/data"

type UserGroupUseCase struct {
	userGroupRepo *data.UserGroupRepo
}

func NewUserGroupUseCase(userGroupRepo *data.UserGroupRepo) *UserGroupUseCase {
	return &UserGroupUseCase{userGroupRepo: userGroupRepo}
}

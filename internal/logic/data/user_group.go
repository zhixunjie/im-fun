package data

type UserGroupRepo struct {
	*Data
}

func NewUserGroupRepo(data *Data) *UserGroupRepo {
	return &UserGroupRepo{
		Data: data,
	}
}

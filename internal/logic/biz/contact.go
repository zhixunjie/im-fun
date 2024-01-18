package biz

import (
	"github.com/zhixunjie/im-fun/internal/logic/data"
)

type ContactUseCase struct {
	repo *data.ContactRepo
}

func NewContactUseCase(repo *data.ContactRepo) *ContactUseCase {
	return &ContactUseCase{repo: repo}
}

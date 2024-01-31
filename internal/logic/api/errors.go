package api

import "errors"

// refer: https://go-kratos.dev/docs/component/errors/

var (
	ErrContactNotExists     = errors.New("contact is not exists")
	ErrSenderTypeNotAllow   = errors.New("sender type not allow")
	ErrReceiverTypeNotAllow = errors.New("sender type not allow")
)

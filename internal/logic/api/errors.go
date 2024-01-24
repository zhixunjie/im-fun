package api

import "errors"

// refer: https://go-kratos.dev/docs/component/errors/

var (
	ErrContactNotExists = errors.New("contact is not exists")
)

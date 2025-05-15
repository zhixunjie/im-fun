package api

import "errors"

// refer: https://go-kratos.dev/docs/component/errors/

var (
	ErrContactNotExists         = errors.New("contact is not exists")
	ErrSenderTypeNotAllow       = errors.New("sender type not allow")
	ErrReceiverTypeNotAllow     = errors.New("receiver type not allow")
	ErrSenderOrReceiverNotAllow = errors.New("sender/receiver not allow")
	ErrMessageBodyNotAllow      = errors.New("message body not allow")
	ErrMessageBodyDecodedFailed = errors.New("message body is decoded error")
	ErrMessageTypeNotAllowed    = errors.New("message type is not allowed")
	ErrMessageContentNotAllowed = errors.New("message content is not allowed")
)

package service

import "errors"

var (
	ErrorInvalidToken       = errors.New("invalid token")
	ErrorNotOnTable         = errors.New("not on table")
	ErrorBadRequest         = errors.New("seat not empty")
	ErrorBadFrame           = errors.New("bad frame")
	ErrorInvalidProtocol    = errors.New("invalid protocol")
	ErrorClientDisconnected = errors.New("client end disconnected")
)

type ErrCode int32

var (
	//客户端错误
	BadRequest      ErrCode = 400
	Unauthorized    ErrCode = 401
	Forbidden       ErrCode = 403
	BadFrame        ErrCode = 460
	InvalidProtocol ErrCode = 461
	//服务器错误
	InternalServerError ErrCode = 500

	UnknowError ErrCode = 521
)

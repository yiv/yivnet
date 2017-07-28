package center

import (
	"errors"
)

var (
	ErrUserExist             = errors.New("user exist")
	ErrNotFound              = errors.New("not found")
	ErrAmountMustBigThanZero = errors.New("amount must big than zero")
	ErrCoinNotEnough         = errors.New("coin not enough")
	ErrGemNotEnough          = errors.New("gem not enough")
	ErrGiftNotEnough         = errors.New("gift not enough")
	ErrBankNotEnough         = errors.New("bank balance not enough")
	ErrPwdEmpty              = errors.New("password can not be empty")
	ErrPwdNotSet             = errors.New("password not set yet")
	ErrPwdWrong              = errors.New("wrong password")
	ErrNickEmpty             = errors.New("nick can not be empty")
	ErrSetValueExist         = errors.New("can not use the same value")
)

package center

import (
	"sync"
)

type User struct {
	Account
	AccountInfo
	sync.RWMutex
	MQ MQRepository
}

type DBRepository interface {
	FindUserById(id int64) (*User, error)
}
type MQRepository interface {
	SaveUser(user *User) (err error)
}

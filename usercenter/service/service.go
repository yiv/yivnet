package service

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"

	"github.com/yiv/yivgame/usercenter/center"
)

type Middleware func(Service) Service

type Service interface {
	GetUserInfo(ctx context.Context, id int64) (user *center.User, err error)
}

var (
	ErrBadRequest = errors.New("bad request")
)

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

type basicService struct {
	DBRepository center.DBRepository
	MQRepository center.MQRepository
	UserManager  *center.UserManager
	logger       log.Logger
}

func NewBasicService(dbRep center.DBRepository, mqRep center.MQRepository, logger log.Logger) Service {
	service := basicService{}
	service.DBRepository = dbRep
	service.MQRepository = mqRep
	service.UserManager = center.NewUserManager(dbRep, mqRep)
	service.logger = logger
	return service
}

//loadUserFromDBtoMem 从数据库中将帐号信息加载到内存中
func (s basicService) loadUserFromDBtoMem(_ context.Context, id int64) (err error) {
	if isExist := s.UserManager.IsUserExist(id); isExist {
		return center.ErrUserExist
	}
	user, err := s.DBRepository.FindUserById(id)
	if err != nil {
		return err
	}
	err = s.UserManager.AddUser(id, user, s.MQRepository)
	if err != nil {
		return err
	}
	return nil
}

//GetUserInfo 获取帐号详细信息
func (s basicService) GetUserInfo(_ context.Context, id int64) (user *center.User, err error) {
	if isExist := s.UserManager.IsUserExist(id); !isExist {
		if err = s.loadUserFromDBtoMem(nil, id); err != nil {
			return nil, err
		}
	}
	return s.UserManager.GetUser(id)
}

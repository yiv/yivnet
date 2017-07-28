package center

type UserMap map[int64]*User

type UserManager struct {
	users        UserMap
	dbRepository DBRepository
	mqRepository MQRepository
}

func NewUserManager(dbRepository DBRepository, mqRepository MQRepository) *UserManager {
	um := &UserManager{}
	umap := make(UserMap)
	um.users = umap
	um.dbRepository = dbRepository
	um.mqRepository = mqRepository
	return um
}

func (um *UserManager) AddUser(uid int64, user *User, mq MQRepository) error {
	if _, ok := um.users[uid]; ok {
		return ErrUserExist
	}
	user.MQ = mq
	um.users[uid] = user
	return nil
}

func (um *UserManager) RemoveUser(uid int64) error {
	return nil
}

func (um *UserManager) GetUser(uid int64) (user *User, err error) {
	user, ok := um.users[uid]
	if ok {
		return user, nil
	}
	user, err = um.dbRepository.FindUserById(uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}
	err = um.AddUser(uid, user, um.mqRepository)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (um *UserManager) IsUserExist(uid int64) bool {
	if _, ok := um.users[uid]; ok {
		return true
	}
	return false
}

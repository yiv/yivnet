package cockroach

import (
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/yiv/yivgame/usercenter/center"
)

const (
	DefalutCoins int64 = 100000
	DefaultGems  int32 = 500
)

type DbRepo struct {
	Conn   *sqlx.DB
	logger log.Logger
}

func NewDbUserRepo(dataSourceName string, logger log.Logger) (center.DBRepository, error) {
	conn, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	dbRepo := new(DbRepo)
	dbRepo.Conn = conn
	dbRepo.logger = logger
	return dbRepo, nil
}

func (repo *DbRepo) FindUserById(uid int64) (user *center.User, err error) {
	user = new(center.User)
	err = repo.Conn.Get(user, "select accounts.id,accounts.uid, accounts.unionid, accounts.uuid , "+
		"accounts.username, accounts.password, accounts.nick, accounts.gender, accounts.addr , accounts.avatar, "+
		"accounts.isguest, accounts.condays, accounts.signdate, accounts.vipsigndate ,accounts.status, accounts.mtime, "+
		"accounts.ctime, accounts.bankpwd , accounts.forbid, accounts.imsi, accounts.imei, accounts.mac, "+
		"accounts.did , accounts.psystem, accounts.pmodel, account_info.coin, account_info.gem, account_info.bank, "+
		"account_info.growth, account_info.level,account_info.viptype, account_info.vipexpiry,account_info.voucher, account_info.token,"+
		"account_info.online from accounts inner join account_info on accounts.uid = account_info.accounts_uid "+
		"where accounts.uid = $1", uid)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error(), "uid", uid)
		return nil, err
	}
	type StringField struct {
		Others  string
		Props   string
		Gifts   string
		Friends string
		Records string
		Tags    string
		Medals  string
	}
	s := new(StringField)
	err = repo.Conn.Get(s, "select accounts.others,account_info.props,account_info.gifts, account_info.friends, account_info.records, "+
		"account_info.tags,account_info.medals from accounts inner join account_info on accounts.uid = account_info.accounts_uid where accounts.uid= $1", uid)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error(), "uid", uid)
		return nil, err
	}
	err = json.Unmarshal([]byte(s.Others), &user.Others)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error(), "Others", s.Others, "msg", "Unmarshal Others")
		return nil, fmt.Errorf("Props json decode err : %s", err.Error())
	}
	err = json.Unmarshal([]byte(s.Props), &user.Props)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error(), "Props", s.Props, "msg", "Unmarshal Props")
		return nil, fmt.Errorf("Props json decode err : %s", err.Error())
	}
	err = json.Unmarshal([]byte(s.Gifts), &user.Gifts)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error(), "Gifts", s.Gifts, "msg", "Unmarshal Gifts")
		return nil, fmt.Errorf("Gifts json decode err : %s", err.Error())
	}
	err = json.Unmarshal([]byte(s.Friends), &user.Friends)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error(), "Friends", s.Friends, "msg", "Unmarshal Friends")
		return nil, fmt.Errorf("Friends json decode err : %s", err.Error())
	}
	err = json.Unmarshal([]byte(s.Records), &user.Records)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error(), "Records", s.Records, "msg", "Unmarshal Records")
		return nil, fmt.Errorf("Records json decode err : %s", err.Error())
	}
	err = json.Unmarshal([]byte(s.Tags), &user.Tags)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error(), "Tags", s.Tags, "msg", "Unmarshal Tags")
		return nil, fmt.Errorf("Tags json decode err : %s", err.Error())
	}
	err = json.Unmarshal([]byte(s.Medals), &user.Medals)
	if err != nil {
		level.Error(repo.logger).Log("err", err.Error(), "Medals", s.Medals, "msg", "Unmarshal Medals")
		return nil, fmt.Errorf("Medals json decode err : %s", err.Error())
	}
	return user, nil
}

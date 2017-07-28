package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/gogo/protobuf/proto"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/yiv/yivgame/usercenter/pb"
	"github.com/yiv/yivgame/usercenter/center"
)

const (
	UserCenterPersistenceTopic string = "UserCenterPersistence"
)

type MQRepo struct {
	asyncProducer sarama.AsyncProducer
	logger        log.Logger
}

func NewMQRepo(addrs []string, logger log.Logger) (center.MQRepository, error) {

	mqRepo := new(MQRepo)
	mqRepo.logger = logger
	producer, err := sarama.NewAsyncProducer(addrs, nil)
	if err != nil {
		level.Error(logger).Log("err", err.Error(), "msg", "error occur on create new kafka producer")
		return nil, err
	}
	mqRepo.asyncProducer = producer
	return mqRepo, nil
}
func (r *MQRepo) SaveUser(user *center.User) (err error) {
	u := &pb.UserInfo{
		Uid:         user.Uid,
		Unionid:     user.Unionid,
		Uuid:        user.Uuid,
		Username:    user.Username,
		Password:    user.Password,
		Nick:        user.Nick,
		Gender:      user.Gender,
		Addr:        user.Addr,
		Avatar:      user.Avatar,
		Isguest:     user.Isguest,
		Condays:     user.Condays,
		Signdate:    user.Signdate,
		Vipsigndate: user.Vipsigndate,
		Status:      user.Status,
		Mtime:       user.Mtime,
		Ctime:       user.Ctime,
		Token:       user.Token,
		Bankpwd:     user.Bankpwd,
		Forbid:      user.Forbid,
		Imsi:        user.Imsi,
		Imei:        user.Imei,
		Mac:         user.Mac,
		Did:         user.Did,
		Psystem:     user.Psystem,
		Pmodel:      user.Pmodel,
		Others:      user.Others,
		Coin:        user.Coin,
		Gem:         user.Gem,
		Bank:        user.Bank,
		Growth:      user.Growth,
		Level:       user.Level,
		Viptype:     user.Viptype,
		Vipexpiry:   user.Vipexpiry,
		Voucher:     user.Voucher,
		Online:      user.Online,
		Props:       user.Props,
		Gifts:       user.Gifts,
		Medals:      user.Medals,
		Friends:     user.Friends,
		Tags:        user.Tags,
		Records:     user.Records,
	}
	pbBytes, err := proto.Marshal(u)
	if err != nil {
		level.Error(r.logger).Log("err", err.Error(), "msg", "error occur on Marshal message for kafka producer")
		return err
	}
	_ = &sarama.ProducerMessage{
		Topic:     UserCenterPersistenceTopic,
		Partition: 0,
		Key:       nil,
		Value:     sarama.ByteEncoder(pbBytes),
	}
	//r.asyncProducer.Input() <- r.newProducerMessage(int32(u.Uid), pbBytes)
	return
}

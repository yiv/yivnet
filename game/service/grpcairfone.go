package service

import (
	"github.com/gogo/protobuf/proto"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/yiv/yivgame/game/gamer"
	"github.com/yiv/yivgame/game/pb"
)

const (
	CmdEnter    uint32 = 10000 //登入
	CmdSendChat uint32 = 10017 //发聊天消息
)
const (
	InformPlayerSendChat uint32 = 20019 //通知玩家发聊天消息
)

const (
	ReportClientOffline uint32 = 40001
)

type airfone struct {
	logger log.Logger
	stream pb.GameService_StreamServer
}

func NewAirfone(stream pb.GameService_StreamServer, logger log.Logger) (af airfone) {
	af = airfone{
		logger: logger,
		stream: stream,
	}
	return
}

//请求返回
func (a airfone) ReplySendChat() (err error) {
	r := &pb.GeRes{Code: 200}
	err = a.send(CmdSendChat, r)
	a.errCheck(err, "ReplySendChat")

	return
}

//通知
func (a airfone) PlayerSendChat(sid gamer.Seat, mid int32, msg string) (err error) {
	r := &pb.ChatMsg{Sid: int32(sid), Mid: int32(mid), Msg: msg}
	err = a.send(InformPlayerSendChat, r)
	a.errCheck(err, "PlayerSendChat")
	return

}
func (a airfone) send(code uint32, pb proto.Message) (err error) {
	var pbBytes []byte
	if pb != nil {
		pbBytes, err = proto.Marshal(pb)
		if err != nil {
			level.Error(a.logger).Log("code", code, "err", err.Error(), "msg", "airfone Marshal err")
			return err
		}
	}
	f := toframe(code, pbBytes)
	err = a.stream.Send(f)
	if err != nil {
		level.Warn(a.logger).Log("code", code, "err", err.Error(), "msg", "airfone send err")
	}
	return
}

func (a airfone) errCheck(err error, method string) {
	if err != nil {
		level.Warn(a.logger).Log("err", err.Error(), "msg", method)
	}
}

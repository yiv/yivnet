package gamer

import (
	"github.com/go-kit/kit/log"
)

type roboAirfone struct {
	logger log.Logger
	player *Player
	table  *Table
}

func NewRobotAirfone(t *Table, logger log.Logger) *roboAirfone {
	r := &roboAirfone{
		logger: log.With(logger, "robot", "airfone"),
		table:  t,
	}
	return r
}

//请求返回
func (r *roboAirfone) ReplySendChat() error {
	return nil
}

//通知
func (r *roboAirfone) PlayerSendChat(sid Seat, mid int32, msg string) error {
	return nil
}

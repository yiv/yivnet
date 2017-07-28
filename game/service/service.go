package service

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/yiv/yivgame/game/gamer"
)

type Middleware func(GameService) GameService

type GameService interface {
	SendChat(ctx context.Context, id gamer.UserID, mid int32, msg string) (err error)
}

type gameService struct {
	logger         log.Logger
	userCenter     gamer.UserCenter
	playerTableMap map[gamer.UserID]*gamer.Table //用于映身已登陆玩家的匹配房间，方便指令的路由，在玩家退出时清除映身，
	rooms          map[gamer.RoomClass]*gamer.Room
}

func NewGameService(fir gamer.RoomOptions, uc *userCenter, logger log.Logger) GameService {
	g := gameService{
		playerTableMap: make(map[gamer.UserID]*gamer.Table),
		rooms:          make(map[gamer.RoomClass]*gamer.Room),
		userCenter:     uc,
		logger:         logger,
	}
	g.rooms[gamer.FirstClassRoom] = gamer.NewRoom(fir, logger, uc)
	return g
}

func (g gameService) SendChat(ctx context.Context, id gamer.UserID, mid int32, msg string) (err error) {
	if id <= 0 {
		return ErrorBadRequest
	}
	var t *gamer.Table
	if t = g.getPlayerTable(id); t == nil {
		return ErrorNotOnTable
	}
	defer g.recoveryTablePanic(t)
	return t.SendChat(id, mid, msg)
}

func (g gameService) recoveryTablePanic(t *gamer.Table) (err error) {
	if e := recover(); e != nil {
		err = fmt.Errorf("recover table panic err: %v", e)
		level.Error(g.logger).Log("table", fmt.Sprintf("%v", t), "PANIC", err)
		delete(g.rooms[t.RoomClass].Tables, t.TableId)
	}
	return
}
func (g gameService) getPlayerTable(id gamer.UserID) (t *gamer.Table) {
	if tb, ok := g.playerTableMap[id]; !ok {
		//映射表找不到玩家
		return nil
	} else {
		return tb
	}
}

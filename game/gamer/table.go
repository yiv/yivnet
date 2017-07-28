package gamer

import (
	"sync"

	"github.com/go-kit/kit/log/level"
)

type (
	TableId  int //TableId 桌号
	SeatCode int //SeatCode 由游戏代码、房间代码、牌桌代码和座位号构成
)

type UserCenter interface {
	GetInfo(id UserID) (info *PlayerInfo, err error)
}

type TableOptions struct {
	TableId   TableId   //桌号
	RoomClass RoomClass //所属房间等级
	BootBet   int64     //底注金币数
}
type Table struct {
	mtx sync.RWMutex
	//牌桌配置
	TableOptions

	//玩家管理
	Seats   map[Seat]UserID
	Players map[UserID]*Player
	//牌局状态
	Status int
	dealer *CardDealer //发牌机
	//外部依赖
	UCenter UserCenter
	logger  Logger
}

//NewTable 新建牌桌
func NewTable(option TableOptions, logger Logger, userCenter UserCenter) (t *Table) {
	t = &Table{
		TableOptions: option,
		Seats:        make(map[Seat]UserID),
		Players:      make(map[UserID]*Player),
		logger:       logger,
		dealer:       NewCardDealer(),
		UCenter:      userCenter,
	}
	return
}

//AddFriend 邀请加好友
func (t *Table) SendChat(id UserID, mid int32, msg string) (err error) {
	var p *Player
	if p = t.getPlayer(id); p == nil {
		return ErrorNotOnTable
	}
	go p.airfone.ReplySendChat()
	level.Debug(t.logger).Log("SendChat", id)
	for _, pl := range t.Players {
		if pl.UserID != p.UserID {
			go pl.airfone.PlayerSendChat(p.Seat, mid, msg)
		}

	}
	return
}

// getPlayer 通过ID获取牌桌上的玩家
func (t *Table) getPlayer(id UserID) (p *Player) {
	if pl, ok := t.Players[id]; ok {
		p = pl
		return
	}
	return
}

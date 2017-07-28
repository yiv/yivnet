package gamer

import (
	"sync"
)

type RoomClass int

const (
	FirstClassRoom RoomClass = iota + 1
	SecondClassRoom
	ThirdClassRoom
)

type RoomOptions struct {
	//牌桌配置
	RoomClass RoomClass
	BootBet   int64 //底注
}
type Room struct {
	RoomOptions
	Tables     map[TableId]*Table
	mtx        sync.RWMutex
	logger     Logger
	userCenter UserCenter
}

//NewRoom 创建新房间
func NewRoom(op RoomOptions, logger Logger, userCenter UserCenter) (r *Room) {
	r = &Room{
		RoomOptions: op,
		logger:      logger,
		userCenter:  userCenter,
		Tables:      make(map[TableId]*Table),
	}
	return
}

//AddTable 房间加桌，当新玩家进入，所有桌都无空座时，需要加桌
func (r *Room) addTable(logger Logger, userCenter UserCenter) (t *Table, err error) {
	if r.CountTable() >= 1000 {
		return nil, ErrorRoomTableExceed
	}
	id := TableId(r.CountTable() + 1)
	tbOption := TableOptions{
		TableId:   id,
		RoomClass: r.RoomClass,
		BootBet:   r.BootBet,
	}
	t = NewTable(tbOption, logger, userCenter)
	r.Tables[id] = t
	return
}

//CountTable 计算开桌局数
func (r *Room) CountTable() int {
	return len(r.Tables)
}

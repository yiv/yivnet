package gamer

type Airfone interface {
	//请求返回
	ReplySendChat() (err error)
	//通知
	PlayerSendChat(sid Seat, mid int32, msg string) (err error)
}

type UserID int64
type Player struct {
	UserID  UserID //玩家个人唯一ID
	Cards   Cards  //玩家牌
	Seat    Seat   //玩家牌桌内的座位号
	IsRobot bool

	//状态
	Status  int //玩家状态
	offline bool

	//传声筒，用于向玩家发送数据
	airfone Airfone

	//玩家个人数据
	Info *PlayerInfo
}
type PlayerInfo struct {
	Token     string
	SeatCode  SeatCode //玩家世界坐标
	Coin      int64
	Gem       int32
	Nick      string
	Avatar    string
	Friends   []UserID
	Character int32
}

func NewPlayer(id UserID, airfone Airfone, info *PlayerInfo, isRobot bool) (player *Player) {
	player = &Player{UserID: id, airfone: airfone, Info: info, IsRobot: isRobot}
	return
}

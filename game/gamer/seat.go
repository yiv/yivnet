package gamer

type (
	Seat int //Seat 玩家在牌桌上的座位号
)

func (s Seat) Next() (seat Seat) {
	if s < 5 {
		return Seat(s + 1)
	}
	return Seat1
}

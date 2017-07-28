package gamer

import "errors"

var (
	ErrorRoomTableExceed = errors.New("tables of room exceed the limit")
	ErrorNotOnTable      = errors.New("not on table")
)

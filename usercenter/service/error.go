package service

import "errors"

var (
	ErrLoadToMem = errors.New("err occur on loading data from db to memory")
)

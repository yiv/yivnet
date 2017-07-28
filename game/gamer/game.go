package gamer

type Logger interface {
	Log(keyvals ...interface{}) error
}

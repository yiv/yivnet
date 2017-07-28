package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/yiv/yivgame/game/gamer"
)

func ServiceLoggingMiddleware(logger log.Logger) Middleware {
	return func(next GameService) GameService {
		return serviceLoggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type serviceLoggingMiddleware struct {
	logger log.Logger
	next   GameService
}

func (mw serviceLoggingMiddleware) SendChat(ctx context.Context, id gamer.UserID, mid int32, msg string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "SendChat",
			"id", id,
			"mid", mid,
			"msg", msg,
			"err", err,
		)
	}(time.Now())
	return mw.next.SendChat(ctx, id, mid, msg)
}

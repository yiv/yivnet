package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/yiv/yivgame/usercenter/center"
)

func ServiceLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return serviceLoggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type serviceLoggingMiddleware struct {
	logger log.Logger
	next   Service
}

//GetUserInfo 获取帐号详细信息
func (mw serviceLoggingMiddleware) GetUserInfo(ctx context.Context, id int64) (user *center.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetUserInfo",
			"id", id,
			"err", err,
			"took", time.Since(begin),
			"user", fmt.Sprintf("%v", user),
		)
	}(time.Now())
	return mw.next.GetUserInfo(ctx, id)
}

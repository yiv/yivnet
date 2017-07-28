package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/yiv/yivgame/usercenter/center"
)

type serviceInstrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func ServiceInstrumentingMiddleware(requestCount metrics.Counter, requestLatency metrics.Histogram) Middleware {
	return func(next Service) Service {
		return serviceInstrumentingMiddleware{
			requestCount:   requestCount,
			requestLatency: requestLatency,
			next:           next,
		}
	}
}

//GetUserInfo 获取帐号详细信息
func (mw serviceInstrumentingMiddleware) GetUserInfo(ctx context.Context, id int64) (user *center.User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetUserInfo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.next.GetUserInfo(ctx, id)
}

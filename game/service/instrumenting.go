package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/yiv/yivgame/game/gamer"
)

type serviceInstrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           GameService
}

func ServiceInstrumentingMiddleware(requestCount metrics.Counter, requestLatency metrics.Histogram) Middleware {
	return func(next GameService) GameService {
		return serviceInstrumentingMiddleware{
			requestCount:   requestCount,
			requestLatency: requestLatency,
			next:           next,
		}
	}
}

func (mw serviceInstrumentingMiddleware) SendChat(ctx context.Context, id gamer.UserID, mid int32, msg string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SendChat", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.next.SendChat(ctx, id, mid, msg)
}

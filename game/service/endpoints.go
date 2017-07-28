package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"

	"github.com/yiv/yivgame/game/gamer"
)

type Endpoints struct {
	Logger           log.Logger
	SendChatEndpoint endpoint.Endpoint
}

type sendChatReq struct {
	Uid int64  `json:"uid"`
	Mid int32  `json:"mid"`
	Msg string `json:"msg"`
}
type sendChatRes struct {
	Err error `json:"err"`
}

func MakeSendChatEndpoint(s GameService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(sendChatReq)
		err = s.SendChat(ctx, gamer.UserID(req.Uid), req.Mid, req.Msg)
		return sendChatRes{Err: err}, err
	}
}

func EndpointInstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)

		}
	}
}

// EndpointLoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func EndpointLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log("error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)

		}
	}
}

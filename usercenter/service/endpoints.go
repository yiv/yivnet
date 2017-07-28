package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"

	"github.com/yiv/yivgame/usercenter/center"
)

type Endpoints struct {
	GetUserInfoEndpoint endpoint.Endpoint
}

type getUserInfoRes struct {
	Uid         int64
	Unionid     string
	Uuid        string
	Username    string
	Password    string
	Nick        string
	Gender      bool
	Addr        string
	Avatar      string
	Isguest     bool
	Condays     int32
	Signdate    int64
	Vipsigndate int64
	Status      bool
	Mtime       int64
	Ctime       int64
	Token       string
	Bankpwd     string
	Forbid      string
	Imsi        string
	Imei        string
	Mac         string
	Did         string
	Psystem     string
	Pmodel      string
	Others      map[string]int32
	Coin        int64
	Gem         int32
	Bank        int64
	Growth      int32
	Level       int32
	Viptype     int32
	Vipexpiry   int64
	Voucher     int32
	Online      int32
	Props       map[string]int32
	Gifts       map[string]int32
	Medals      map[string]int32
	Friends     []int64
	Tags        []string
	Records     map[string]int32
	Err         error `json:"err"`
}

func MakeGetUserInfoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		id := request.(int64)
		u, err := s.GetUserInfo(ctx, id)
		if err != nil {
			return nil, err
		}
		return getUserInfoRes{
			Uid:         u.Uid,
			Unionid:     u.Unionid,
			Uuid:        u.Uuid,
			Username:    u.Username,
			Password:    u.Password,
			Nick:        u.Nick,
			Gender:      u.Gender,
			Addr:        u.Addr,
			Avatar:      u.Avatar,
			Isguest:     u.Isguest,
			Condays:     u.Condays,
			Signdate:    u.Signdate,
			Vipsigndate: u.Vipsigndate,
			Status:      u.Status,
			Mtime:       u.Mtime,
			Ctime:       u.Ctime,
			Token:       u.Token,
			Bankpwd:     u.Bankpwd,
			Forbid:      u.Forbid,
			Imsi:        u.Imsi,
			Imei:        u.Imei,
			Mac:         u.Mac,
			Did:         u.Did,
			Psystem:     u.Psystem,
			Pmodel:      u.Pmodel,
			Others:      u.Others,
			Coin:        u.Coin,
			Gem:         u.Gem,
			Bank:        u.Bank,
			Growth:      u.Growth,
			Level:       u.Level,
			Viptype:     u.Viptype,
			Vipexpiry:   u.Vipexpiry,
			Voucher:     u.Voucher,
			Online:      u.Online,
			Props:       u.Props,
			Gifts:       u.Gifts,
			Medals:      u.Medals,
			Friends:     u.Friends,
			Tags:        u.Tags,
			Records:     u.Records,
		}, err
	}
}
func (e Endpoints) GetUserInfo(ctx context.Context, uid int64) (user *center.User, err error) {
	response, err := e.GetUserInfoEndpoint(ctx, uid)
	if err != nil {
		return nil, err
	}
	r := response.(getUserInfoRes)
	return &center.User{
		Account: center.Account{
			Uid:         r.Uid,
			Unionid:     r.Unionid,
			Uuid:        r.Uuid,
			Username:    r.Username,
			Password:    r.Password,
			Nick:        r.Nick,
			Gender:      r.Gender,
			Addr:        r.Addr,
			Avatar:      r.Avatar,
			Isguest:     r.Isguest,
			Condays:     r.Condays,
			Signdate:    r.Signdate,
			Vipsigndate: r.Vipsigndate,
			Status:      r.Status,
			Mtime:       r.Mtime,
			Ctime:       r.Ctime,
			Bankpwd:     r.Bankpwd,
			Forbid:      r.Forbid,
			Imsi:        r.Imsi,
			Imei:        r.Imei,
			Mac:         r.Mac,
			Did:         r.Did,
			Psystem:     r.Psystem,
			Pmodel:      r.Pmodel,
			Others:      r.Others,
		},
		AccountInfo: center.AccountInfo{
			Token:     r.Token,
			Coin:      r.Coin,
			Gem:       r.Gem,
			Bank:      r.Bank,
			Growth:    r.Growth,
			Level:     r.Level,
			Viptype:   r.Viptype,
			Vipexpiry: r.Vipexpiry,
			Voucher:   r.Voucher,
			Online:    r.Online,
			Props:     r.Props,
			Gifts:     r.Gifts,
			Medals:    r.Medals,
			Friends:   r.Friends,
			Tags:      r.Tags,
			Records:   r.Records,
		},
	}, nil
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

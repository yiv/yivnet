package service

import (
	"context"
	"io"
	"time"

	stdopentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/sd"
	ketcd "github.com/go-kit/kit/sd/etcd"
	"github.com/go-kit/kit/sd/lb"

	"fmt"
	"github.com/yiv/yivgame/game/gamer"
	ucgrpccli "github.com/yiv/yivgame/usercenter/client"
	"github.com/yiv/yivgame/usercenter/service"
)

type userCenter struct {
	logger    log.Logger
	endpoints service.Endpoints
}

func (u *userCenter) GetInfo(id gamer.UserID) (info *gamer.PlayerInfo, err error) {
	var ctx = context.Background()
	user, err := u.endpoints.GetUserInfo(ctx, int64(id))
	if err != nil {
		level.Error(u.logger).Log("userCenter", "GetInfo", "id", id, "err", err.Error())
		return nil, err
	}
	var friends []gamer.UserID
	for _, f := range user.Friends {
		friends = append(friends, gamer.UserID(f))
	}
	info = &gamer.PlayerInfo{
		Token:     user.Token,
		SeatCode:  gamer.SeatCode(user.Online),
		Coin:      user.Coin,
		Gem:       user.Gem,
		Nick:      user.Nick,
		Avatar:    user.Avatar,
		Friends:   friends,
		Character: user.Others["character"],
	}
	return
}

func NewUserCenter(serviceName string, etcdAddr []string, retryMax int, retryTimeout time.Duration, logger log.Logger) (uc *userCenter, err error) {
	var ctx = context.Background()
	tracer := stdopentracing.GlobalTracer() // no-op
	client, err := ketcd.NewClient(ctx, etcdAddr, ketcd.ClientOptions{})
	if err != nil {
		level.Error(logger).Log("userCenter", "NewUserCenter", "err", err.Error())
		return nil, fmt.Errorf("NewClient err : %s", err.Error())
	}
	instancer, err := ketcd.NewInstancer(client, serviceName, logger)

	endpoints := service.Endpoints{}
	{
		factory := factoryFor(service.MakeGetUserInfoEndpoint, tracer, logger)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.GetUserInfoEndpoint = retry
	}
	uc = &userCenter{
		logger:    logger,
		endpoints: endpoints,
	}
	return
}
func factoryFor(makeEndpoint func(service.Service) endpoint.Endpoint, tracer stdopentracing.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			level.Error(logger).Log("userCenter", "factoryFor", "err", err.Error())
			return nil, nil, err
		}
		svr := ucgrpccli.New(conn, tracer, logger)
		ep := makeEndpoint(svr)

		return ep, conn, nil
	}
}

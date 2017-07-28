package service

import (
	"context"

	stdopentracing "github.com/opentracing/opentracing-go"
	oldcontext "golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"

	"github.com/yiv/yivgame/usercenter/pb"
)

type grpcHandler struct {
	getUserInfo kitgrpc.Handler
}

func MakeGRPCHandler(endpoints Endpoints, tracer stdopentracing.Tracer, logger log.Logger) pb.UserServer {
	options := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorLogger(logger),
	}
	return &grpcHandler{
		getUserInfo: kitgrpc.NewServer(
			endpoints.GetUserInfoEndpoint,
			DecodeGRPCGetUserInfoReq,
			EncodeGRPCGetUserInfoRes,
			append(options, kitgrpc.ServerBefore(opentracing.GRPCToContext(tracer, "getUserInfo", logger)))...,
		),
	}

}

//for server
func DecodeGRPCGetUserInfoReq(_ context.Context, request interface{}) (interface{}, error) {
	r := request.(*pb.UserId)
	return r.Uid, nil
}
func EncodeGRPCGetUserInfoRes(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(getUserInfoRes)
	return &pb.UserInfo{
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
		Token:       r.Token,
		Bankpwd:     r.Bankpwd,
		Forbid:      r.Forbid,
		Imsi:        r.Imsi,
		Imei:        r.Imei,
		Mac:         r.Mac,
		Did:         r.Did,
		Psystem:     r.Psystem,
		Pmodel:      r.Pmodel,
		Others:      r.Others,
		Coin:        r.Coin,
		Gem:         r.Gem,
		Bank:        r.Bank,
		Growth:      r.Growth,
		Level:       r.Level,
		Viptype:     r.Viptype,
		Vipexpiry:   r.Vipexpiry,
		Voucher:     r.Voucher,
		Online:      r.Online,
		Props:       r.Props,
		Gifts:       r.Gifts,
		Medals:      r.Medals,
		Friends:     r.Friends,
		Tags:        r.Tags,
		Records:     r.Records,
		Err:         err2str(r.Err),
	}, nil
}

//for client
func EncodeGRPCGetUserInfoReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(int64)
	return &pb.UserId{Uid: req}, nil
}
func DecodeGRPCGetUserInfoRes(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(*pb.UserInfo)
	return getUserInfoRes{
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
		Token:       r.Token,
		Bankpwd:     r.Bankpwd,
		Forbid:      r.Forbid,
		Imsi:        r.Imsi,
		Imei:        r.Imei,
		Mac:         r.Mac,
		Did:         r.Did,
		Psystem:     r.Psystem,
		Pmodel:      r.Pmodel,
		Others:      r.Others,
		Coin:        r.Coin,
		Gem:         r.Gem,
		Bank:        r.Bank,
		Growth:      r.Growth,
		Level:       r.Level,
		Viptype:     r.Viptype,
		Vipexpiry:   r.Vipexpiry,
		Voucher:     r.Voucher,
		Online:      r.Online,
		Props:       r.Props,
		Gifts:       r.Gifts,
		Medals:      r.Medals,
		Friends:     r.Friends,
		Tags:        r.Tags,
		Records:     r.Records,
		Err:         str2err(r.Err),
	}, nil
}
func (g *grpcHandler) GetUserInfo(ctx oldcontext.Context, req *pb.UserId) (*pb.UserInfo, error) {
	_, rep, err := g.getUserInfo.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserInfo), nil
}

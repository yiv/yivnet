// Package grpc provides a gRPC client for the add service.
package grpc

import (
	_ "time"

	jujuratelimit "github.com/juju/ratelimit"
	stdopentracing "github.com/opentracing/opentracing-go"
	//"github.com/sony/gobreaker"
	"google.golang.org/grpc"

	//"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	"github.com/yiv/yivgame/usercenter/service"
	"github.com/yiv/yivgame/usercenter/pb"
)

func New(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) service.Service {

	limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))

	var getDeviceIDEndpoint endpoint.Endpoint
	{
		getDeviceIDEndpoint = grpctransport.NewClient(
			conn,
			"pb.User",
			"GetDeviceID",
			service.EncodeGRPCGetDeviceIDReq,
			service.DecodeGRPCGetDeviceIDRes,
			pb.GetDeviceIDRes{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		getDeviceIDEndpoint = opentracing.TraceClient(tracer, "GetDeviceID")(getDeviceIDEndpoint)
		getDeviceIDEndpoint = limiter(getDeviceIDEndpoint)

	}
	var loginGuestEndpoint endpoint.Endpoint
	{
		loginGuestEndpoint = grpctransport.NewClient(
			conn,
			"pb.User",
			"LoginGuest",
			service.EncodeGRPCLoginGuestReq,
			service.DecodeGRPCLoginGuestRes,
			pb.LoginRes{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		loginGuestEndpoint = opentracing.TraceClient(tracer, "LoginGuest")(loginGuestEndpoint)
		loginGuestEndpoint = limiter(loginGuestEndpoint)

	}
	var getUserInfoEndpoint endpoint.Endpoint
	{
		getUserInfoEndpoint = grpctransport.NewClient(
			conn,
			"pb.User",
			"GetUserInfo",
			service.EncodeGRPCGetUserInfoReq,
			service.DecodeGRPCGetUserInfoRes,
			pb.UserInfo{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		getUserInfoEndpoint = opentracing.TraceClient(tracer, "GetUserInfo")(getUserInfoEndpoint)
		getUserInfoEndpoint = limiter(getUserInfoEndpoint)
	}
	var adjustCoinEndpoint endpoint.Endpoint
	{
		adjustCoinEndpoint = grpctransport.NewClient(
			conn,
			"pb.User",
			"AdjustCoin",
			service.EncodeGRPCAdjustCoinReq,
			service.DecodeGRPCAdjustCoinRes,
			pb.AdjustCoinRes{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		adjustCoinEndpoint = opentracing.TraceClient(tracer, "AdjustCoin")(adjustCoinEndpoint)
		adjustCoinEndpoint = limiter(adjustCoinEndpoint)
	}
	var adjustGiftEndpoint endpoint.Endpoint
	{
		adjustGiftEndpoint = grpctransport.NewClient(
			conn,
			"pb.User",
			"AdjustGift",
			service.EncodeGRPCAdjustGiftReq,
			service.DecodeGRPCAdjustGiftRes,
			pb.AdjustGiftRes{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		adjustGiftEndpoint = opentracing.TraceClient(tracer, "AdjustGift")(adjustGiftEndpoint)
		adjustGiftEndpoint = limiter(adjustGiftEndpoint)
	}
	var adjustGemEndpoint endpoint.Endpoint
	{
		adjustGemEndpoint = grpctransport.NewClient(
			conn,
			"pb.User",
			"AdjustGem",
			service.EncodeGRPCAdjustGemReq,
			service.DecodeGRPCAdjustGemRes,
			pb.AdjustGemRes{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		adjustGemEndpoint = opentracing.TraceClient(tracer, "AdjustGem")(adjustGemEndpoint)
		adjustGemEndpoint = limiter(adjustGemEndpoint)
	}
	var updateSeatCodeEndpoint endpoint.Endpoint
	{
		updateSeatCodeEndpoint = grpctransport.NewClient(
			conn,
			"pb.User",
			"UpdateSeatCode",
			service.EncodeGRPCUpdateSeatCodeReq,
			service.DecodeGRPCUpdateSeatCodeRes,
			pb.UpdateSeatCodeRes{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		updateSeatCodeEndpoint = opentracing.TraceClient(tracer, "UpdateSeatCode")(updateSeatCodeEndpoint)
		updateSeatCodeEndpoint = limiter(updateSeatCodeEndpoint)
	}
	return service.Endpoints{
		GetDeviceIDEndpoint:    getDeviceIDEndpoint,
		LoginGuestEndpoint:     loginGuestEndpoint,
		GetUserInfoEndpoint:    getUserInfoEndpoint,
		AdjustCoinEndpoint:     adjustCoinEndpoint,
		AdjustGiftEndpoint:     adjustGiftEndpoint,
		AdjustGemEndpoint:      adjustGemEndpoint,
		UpdateSeatCodeEndpoint: updateSeatCodeEndpoint,
	}

}

//func makeEndpoint(serviceName, method string, conn *grpc.ClientConn, grpcReply interface{}, tracer stdopentracing.Tracer, logger log.Logger, limiter endpoint.Middleware, enc grpctransport.EncodeRequestFunc, dec grpctransport.DecodeResponseFunc) endpoint.Endpoint {
//	var endpoint endpoint.Endpoint
//	endpoint = grpctransport.NewClient(
//		conn,
//		serviceName,
//		method,
//		enc,
//		dec,
//		grpcReply,
//		grpctransport.ClientBefore(opentracing.ToGRPCRequest(tracer, logger)),
//	).Endpoint()
//	endpoint = opentracing.TraceClient(tracer, method)(endpoint)
//	endpoint = limiter(endpoint)
//	//endpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
//	//Name:    method,
//	//Timeout: 1 * time.Second,
//	//}))(endpoint)
//	return endpoint
//}

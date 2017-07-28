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

	"github.com/yiv/yivgame/usercenter/pb"
	"github.com/yiv/yivgame/usercenter/service"
)

func New(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) service.Service {

	limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))

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

	return service.Endpoints{
		GetUserInfoEndpoint: getUserInfoEndpoint,
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

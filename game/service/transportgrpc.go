package service

import (
	"context"
	"io"

	"github.com/gogo/protobuf/proto"
	stdopentracing "github.com/opentracing/opentracing-go"
	oldcontext "golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"

	"github.com/yiv/yivgame/game/pb"
)

type grpcHandler struct {
	logger   log.Logger
	sendChat kitgrpc.Handler
}

func MakeGRPCHandler(endpoints Endpoints, tracer stdopentracing.Tracer, logger log.Logger) pb.GameServiceServer {
	options := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorLogger(logger),
	}
	return &grpcHandler{
		logger: logger,
		sendChat: kitgrpc.NewServer(
			endpoints.SendChatEndpoint,
			DecodeGRPCSendChatReq,
			EncodeGRPCSendChatRes,
			append(options, kitgrpc.ServerBefore(opentracing.GRPCToContext(tracer, "sendChat", logger)))...,
		),
	}

}

//sendChat
func DecodeGRPCSendChatReq(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(sendChatReq)
	return req, nil
}
func EncodeGRPCSendChatRes(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(sendChatRes)
	return nil, res.Err
}

//Stream 实现grpc双向流的stream方法
func (g *grpcHandler) Stream(stream pb.GameService_StreamServer) (err error) {
	var userId int64
	dieCH := make(chan struct{})
	defer func() {
		close(dieCH)
	}()
	recvChan := g.goRecv(stream, dieCH)
	for {
		select {
		case frame, ok := <-recvChan: // frames from agent
			if !ok { // EOF
				level.Info(g.logger).Log("FUNC", "Stream", "msg", "stream receive err, exit send loop")
				if err = g.cmdRoute(&userId, ReportClientOffline, nil, stream); err != nil {
					level.Info(g.logger).Log("err", err.Error(), "msg", "err when report client stream disconnected")
					return
				}
				return ErrorClientDisconnected
			}
			code := frameCode(frame)
			req := framePBbytes(frame)
			err = g.cmdRoute(&userId, code, req, stream)
			if err != nil {
				errCode := g.ErrToCode(err)
				f := g.ErrCodeToFrame(code, errCode)
				if err = stream.Send(f); err != nil {
					level.Info(g.logger).Log("code", code, "err", err.Error(), "msg", "stream err on send err response")
					return err
				}
			}
		}
	}
	return nil
}

//为
func (g *grpcHandler) goRecv(stream pb.GameService_StreamServer, dieCH chan struct{}) chan *pb.Frame {
	recvChan := make(chan *pb.Frame, 1)
	go func() {
		defer func() {
			close(recvChan)
		}()
		for {
			frame, err := stream.Recv()
			if err == io.EOF {
				level.Error(g.logger).Log("FUNC", "goRecv", "err", err.Error(), "msg", "stream receive err, exit recv loop")
				return
			}
			if err != nil {
				level.Error(g.logger).Log("FUNC", "goRecv", "err", err.Error(), "msg", "stream receive err, , exit recv loop")
				return
			}
			//level.Debug(g.logger).Log("frame", frame.Payload, "msg", "stream received")
			select {
			case recvChan <- frame:
			case <-dieCH:
			}
		}
	}()
	return recvChan
}

//cmdRoute 对收到命令代码进行路由
func (g *grpcHandler) cmdRoute(userId *int64, code uint32, req []byte, stream pb.GameService_StreamServer) error {
	level.Debug(g.logger).Log("CMD", "cmdRoute", "userId", *userId, "code", code)
	ctx := oldcontext.Background()

	switch code {
	case CmdEnter:
		request := &pb.EnterTableReq{}
		//记录当前登陆的玩家ID，作为对本stream已验证的标记
		*userId = request.Uid

		return nil
		//牌局命令
		if *userId <= 0 {
			//连接没有经过鉴权
			level.Debug(g.logger).Log("CMD", "cmdRoute", "userId", *userId, "code", code, "msg", "stream unauthorized")
			return ErrorInvalidToken
		}

		switch code {
		case CmdSendChat:
			request := &pb.ChatReq{}
			proto.Unmarshal(req, request)
			_, _, err := g.sendChat.ServeGRPC(ctx, sendChatReq{Uid: *userId, Mid: request.Mid, Msg: request.Msg})
			if err != nil {
				level.Error(g.logger).Log("userId", *userId, "err", err.Error(), "msg", "sendGift serve err")
				return err
			}
			return nil
		default:
			return ErrorInvalidProtocol
		}
		return nil
	}
	return ErrorInvalidProtocol
}
func (g *grpcHandler) ErrToCode(err error) ErrCode {
	switch err {
	case ErrorBadFrame:
		return BadFrame
	case ErrorInvalidProtocol:
		return InvalidProtocol
	case ErrorInvalidToken:
		return Unauthorized
	}
	return UnknowError
}

func (g *grpcHandler) ErrCodeToFrame(code uint32, errCode ErrCode) *pb.Frame {
	var pbBytes []byte
	switch code {
	default:
		pbBytes, _ = proto.Marshal(&pb.GeRes{Code: int32(errCode)})
	}
	return toframe(code, pbBytes)
}
func (g *grpcHandler) invalidProtocol(code uint32) *pb.Frame {
	level.Error(g.logger).Log("code", code, "msg", "protocol code invalid")
	pbBytes, _ := proto.Marshal(&pb.GeRes{Code: int32(InvalidProtocol)})
	return toframe(code, pbBytes)
}

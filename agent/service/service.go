package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/yiv/yivgame/game/pb"
)

const (
	WebSocketReadDeadline int = 15 //秒
)

var (
	ErrReadWebSocket  = errors.New("err on read data from webSocket")
	ErrWriteWebSocket = errors.New("err on write data to webSocket")
	ErrReadRPCStream  = errors.New("err on read data from rpc stream")
	ErrWriteRPCStream = errors.New("err on write data to rpc stream")
)

type Session struct {
	id      string //会话随机ID
	logger  log.Logger
	as      *AgentService
	conn    *websocket.Conn             //webSocket连接
	stream  pb.GameService_StreamClient //rpc stream，每会话每stream
	wg      sync.WaitGroup
	dieChan chan struct{}
}

func NewSession(id string, conn *websocket.Conn, stream pb.GameService_StreamClient, service *AgentService, logger log.Logger) *Session {
	s := &Session{
		id:      id,
		as:      service,
		logger:  logger,
		conn:    conn,
		stream:  stream,
		dieChan: make(chan struct{}),
	}
	return s
}

func (s *Session) ForwardToClient(payload []byte) (err error) {
	err = s.conn.WriteMessage(websocket.BinaryMessage, payload)
	level.Debug(s.logger).Log("id", s.id, "msg", "forward server data to client")
	if err != nil {
		level.Error(s.logger).Log("err", err.Error(), "msg", "webSocket write err")
	}
	return
}
func (s *Session) ForwardToServer(f *pb.Frame) (err error) {
	//转发到服务器
	err = s.stream.Send(f)
	level.Debug(s.logger).Log("id", s.id, "msg", "forward client data to server")
	if err != nil {
		level.Error(s.logger).Log("err", err.Error(), "msg", "stream send err")
	}
	return
}
func (s *Session) HandleClient() {
	errChan := make(chan error)
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		//处理来自客户端的数据
		for {
			//s.conn.SetReadDeadline(time.Now().Add(time.Duration(WebSocketReadDeadline) * time.Second))
			_, bytes, err := s.conn.ReadMessage()
			if err != nil {
				level.Error(s.logger).Log("id", s.id, "err", err.Error(), "msg", "session webSocket read data err")
				errChan <- ErrReadWebSocket
				return
			}
			err = s.ForwardToServer(&pb.Frame{Payload: bytes})
			if err != nil {
				errChan <- ErrWriteRPCStream
				return
			}
			select {
			case <-s.dieChan:
				return
			default:
			}

		}
	}()

	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		//处理来自服务器的数据
		for {
			f, err := s.stream.Recv()
			if err != nil {
				level.Error(s.logger).Log("err", err.Error(), "msg", "stream recv err")
				errChan <- ErrReadRPCStream
				return
			}
			err = s.ForwardToClient(f.Payload)
			if err != nil {
				errChan <- ErrWriteWebSocket
				return
			}
			select {
			case <-s.dieChan:
				return
			default:

			}
		}

	}()
	s.as.closeSession(s.id, <-errChan)
}

type AgentService struct {
	mtx              sync.RWMutex
	SessDieChan      chan int64
	Sessions         map[string]*Session
	logger           log.Logger
	gameServerClient pb.GameServiceClient
}

func NewAgentService(gameServerClient pb.GameServiceClient, logger log.Logger) AgentService {
	level.Info(logger).Log("msg", "start new agent service")

	agent := AgentService{
		gameServerClient: gameServerClient,
		logger:           logger,
		Sessions:         make(map[string]*Session),
	}
	return agent
}

func (a AgentService) WebSocketServer(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		level.Error(a.logger).Log("err", err.Error(), "msg", "err on create webSocket for session")
		return
	}
	stream, err := a.gameServerClient.Stream(context.Background())
	if err != nil {
		conn.Close()
		level.Error(a.logger).Log("err", err.Error(), "msg", "err on create rpc stream for session")
		return
	}
	a.mtx.Lock()
	defer a.mtx.Unlock()
	id := uuid.NewV4().String()
	a.Sessions[id] = NewSession(id, conn, stream, &a, a.logger)
	go a.Sessions[id].HandleClient()
	fmt.Println("start new session")
}

func (a AgentService) closeSession(id string, err error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	bf := len(a.Sessions)
	ss := a.Sessions[id]
	ss.conn.Close()
	ss.stream.CloseSend()
	close(ss.dieChan)
	ss.wg.Wait()
	delete(a.Sessions, id)
	af := len(a.Sessions)
	level.Error(a.logger).Log("id", id, "before", bf, "after", af, "err", err.Error(), "msg", "close session")
}

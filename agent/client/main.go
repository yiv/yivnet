package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"

	"github.com/yiv/yivgame/game/pb"
)

var (
	agentAddr string = ":10050"
)

type client struct {
	sk        *websocket.Conn
	uid       int64
	token     string
	roomClass int32
	player    *pb.Player
}

func newClient(uid int64) (c *client) {
	c = &client{
		sk:        newWebSocket(),
		uid:       uid,
		token:     strconv.FormatInt(uid, 10),
		roomClass: 3,
	}

	return
}

func (c *client) recv(errChan chan error) {
	for {
		_, message, err := c.sk.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			errChan <- err
			return
		}
		code, pbBytes := splitBytes(message)
		if code == 10018 {
			res := &pb.GeRes{}
			proto.Unmarshal(pbBytes, res)
			fmt.Printf("[返回]发消息[code: %v]\n", res.Code)
		} else if code == 20019 {
			res := &pb.ChatMsg{}
			proto.Unmarshal(pbBytes, res)
			fmt.Printf("[通知]信息[发者: %v,编号：%v]\n", res.Sid, res.Mid)
		} else {
			res := &pb.GeRes{}
			proto.Unmarshal(pbBytes, res)
			fmt.Printf("[返回]未知[code: %v]\n", res.Code)
		}
		fmt.Printf("\n 17=消息\n\n")
	}
}
func (c *client) send(errChan chan error) {
	var cmdCode uint32
	fmt.Println("=======命令行启动=====")
	for {
		fmt.Scanln(&cmdCode)

		if cmdCode == 100 {
			c.sk.Close()
			continue
		} else if cmdCode == 99 {
			c.enterTable()
			continue
		} else {
			cmdCode += 10000
		}

		fmt.Printf("输入命令: %v\n", cmdCode)
		err := c.sk.WriteMessage(websocket.BinaryMessage, toBytes(cmdCode, nil))
		if err != nil {
			fmt.Println("命令发送错误：", err)
			errChan <- err
			return
		}
	}

}

func newWebSocket() *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: agentAddr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	fmt.Println("socket connected ")
	return c
}
func splitBytes(payload []byte) (code uint32, pbBytes []byte) {
	code = uint32(payload[0])<<24 | uint32(payload[1])<<16 | uint32(payload[2])<<8 | uint32(payload[3])
	pbBytes = payload[4:]
	return
}
func toBytes(i uint32, pbBytes []byte) (payload []byte) {
	payload = append(payload, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	payload = append(payload, pbBytes...)
	return
}

var (
	uid = flag.Int("uid", 1, "the client uid")
)

var (
	userMap = map[int]int64{
		1: 953685341995009,
		2: 2232870895617,
	}
)

func main() {
	flag.Parse()
	uid := userMap[*uid]
	errChan := make(chan error)
	cli := newClient(uid)
	go cli.recv(errChan)
	go cli.send(errChan)
	fmt.Println("terminate by err :", <-errChan)
}

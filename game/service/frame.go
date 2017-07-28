package service

import (
	"github.com/yiv/yivgame/game/pb"
)

func code2bytes(i uint32) (b []byte) {
	b = append(b, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	return
}
func bytes2code(b []byte) (i uint32) {
	i = uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	return
}

func frameCode(f *pb.Frame) uint32 {
	payload := f.Payload
	return bytes2code(payload[0:4])
}
func framePBbytes(f *pb.Frame) []byte {
	payload := f.Payload
	return payload[4:]
}
func toframe(i uint32, pbBytes []byte) *pb.Frame {
	payload := append(code2bytes(i), pbBytes...)
	return &pb.Frame{Payload: payload}
}

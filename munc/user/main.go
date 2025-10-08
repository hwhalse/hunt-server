package user

import (
	"eolian/munc/proto/munc/roomData"
	"net"
)

type User struct {
	Nickname     string
	TCPAddress   net.TCPAddr
	Host         string
	PositionAddr string
	VoipAddr     string
	UserId       string
	Socket       net.TCPConn
	Room         string
	Platform     roomData.Platform
	Connected    bool
}

package voip_server

import (
	handle_errors "eolian/munc/errors"
	"eolian/munc/state"
	"fmt"
	"log/slog"
	"net"
	"os"
	"syscall"
)

var tunnel = make(map[string]net.Addr)
var userRooms = make(map[string]string)

func RemoveFromUserRooms(addr string) {
	delete(userRooms, addr)
}

func handle(ch chan struct {
	buf []byte
	src net.Addr
}, conn net.PacketConn) {
	for frame := range ch {
		f := frame.buf
		room := state.State.GetRoomByVoipAddr(frame.src.String())
		tunnel[frame.src.String()] = frame.src
		if room.Name == "" {
			continue
		}
		userRooms[frame.src.String()] = room.Name
		for k, v := range userRooms {
			if k == frame.src.String() {
				continue
			}
			if v == room.Name {
				addr, err := net.ResolveUDPAddr("udp", k)
				handle_errors.CheckError(err, "check")
				_, err = conn.WriteTo(f, addr)
				handle_errors.CheckError(err, "Unable to write")
				continue
			}
		}
		continue
	}
}

func StartVoipServer() {
	voipPort := os.Getenv("MUNC_VOICE_PORT")
	slog.Debug("Voip Server Started", "port", voipPort)
	fmt.Println("voip port ", voipPort)
	voipAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:"+voipPort)
	if err != nil {
		fmt.Println("Unable to bind udp server to port "+voipPort, err)
		return
	}

	udpSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)

	handle_errors.CheckError(err, "Unable to create dgram socket")

	err = syscall.SetsockoptInt(udpSocket, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	handle_errors.CheckError(err, "Unable to set so reuseaddr")

	err = syscall.Bind(udpSocket, &syscall.SockaddrInet4{Port: voipAddr.Port})

	handle_errors.CheckError(err, "Cannot bind socket")

	file := os.NewFile(uintptr(udpSocket), string(rune(udpSocket)))
	conn, err := net.FilePacketConn(file)
	handle_errors.CheckError(err, "Cannot create connection from socket")

	in := make(chan struct {
		buf []byte
		src net.Addr
	}, 16)

	go handle(in, conn)

	for {
		buf := make([]byte, 1600)
		n, src, err := conn.ReadFrom(buf)
		handle_errors.CheckError(err, "unable to read")
		in <- struct {
			buf []byte
			src net.Addr
		}{buf: buf[:n], src: src}
	}
}

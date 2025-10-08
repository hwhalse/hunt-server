package types

import (
	"net"
	"sync"
)

type OutChan struct {
	Channel chan OutboundMessage
}

type OutboundMessage struct {
	Bytes []byte
	Conn  net.TCPConn
}

type Connections struct {
	sync.Mutex
	ConnMap map[string]*net.TCPConn
}

package tcp_server

import (
	"eolian/munc/tcp_server/types"
	"github.com/rs/zerolog/log"
	"net"
)

func StartTCPServer() error {
	//initialize state
	outChannel := types.OutChan{Channel: make(chan types.OutboundMessage, 100)}
	connectionMap := types.Connections{ConnMap: make(map[string]*net.TCPConn)}
	rooms := types.Rooms{RoomMap: map[string]*types.Room{}}

	go sendTcpMessage(&outChannel)
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: 5000,
	})
	if err != nil {
		log.Error().Err(err).Msg("Unable to start TCP Server")
		return err
	}
	log.Info().Msg("TCP Server started on port 5000.")
	defer listener.Close()
	for {
		conn, err := listener.AcceptTCP()
		connectionMap.Lock()
		connectionMap.ConnMap[conn.RemoteAddr().String()] = conn
		log.Debug().Str("remote address", conn.RemoteAddr().String()).Msg("New TCP Connection")
		if err != nil {
			log.Error().Err(err).Msg("Unable to accept incoming TCP conn")
			return err
		}
		go readLoop(conn, &connectionMap, &outChannel, &rooms)
	}
}

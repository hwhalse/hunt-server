package tcp_server

import (
	"encoding/binary"
	"eolian/munc/constants"
	"eolian/munc/proto/munc/core"
	"eolian/munc/tcp_server/handleEvent"
	"eolian/munc/tcp_server/types"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"net"
)

func readLoop(
	conn *net.TCPConn,
	conns *types.Connections,
	outChan *types.OutChan,
	rooms *types.Rooms,
) {
	defer conn.Close()
	for {
		lenBytes := make([]byte, 4)
		_, err := conn.Read(lenBytes)
		if err != nil {
			log.Error().Err(err).Msg("readLoop error reading length")
			if err.Error() == "EOF" {
				log.Error().Str("remote address", conn.RemoteAddr().String()).Msg("Connection EOF")
				break
			}
		}
		length := binary.BigEndian.Uint32(lenBytes)
		protoBytes := make([]byte, length)
		_, err = conn.Read(protoBytes)
		muncMessage := &core.PbNetworkReq{}
		err = proto.Unmarshal(protoBytes, muncMessage)
		if err != nil {
			log.Error().Err(err).Msg("readLoop error unmarshaling protobuf")
		}
		muncEventId := muncMessage.MuncEventId
		log.Debug().Uint32("eventId", muncEventId).Uint32("length", length).Msg("Incoming TCP Message")
		switch muncEventId {
		case constants.TCP_INIT:
			err = handleTCPInit(muncMessage, *conn, outChan)
			if err != nil {
				log.Error().Err(err).Msg("Handle TCP Init error")
				return
			}
		case constants.TCP_JOIN_ROOM:
			err = handleRoomJoin(muncMessage, *conn, rooms, outChan)
			if err != nil {
				log.Error().Err(err).Msg("Handle Room Join error")
				return
			}
		case constants.TCP_EXIT_ROOM:
			handleRoomExit(muncMessage, *conn, rooms, *outChan)
		case constants.GET_LOCATIONS:
			err = handleEvent.GetLocations(muncMessage, *conn, rooms, outChan)
			if err != nil {
				log.Error().Err(err).Msg("Handle Event GetLocations error")
				return
			}
		case constants.UPDATE_LOCATION:
			err = handleEvent.UpdateLocation(muncMessage, *conn, rooms, outChan)
			if err != nil {
				log.Error().Err(err).Msg("Handle Event UpdateLocation error")
				return
			}
		}
	}
}

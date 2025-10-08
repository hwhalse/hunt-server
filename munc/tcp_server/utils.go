package tcp_server

import (
	"encoding/binary"
	"eolian/munc/constants"
	handle_errors "eolian/munc/errors"
	"eolian/munc/proto/munc/core"
	"eolian/munc/proto/munc/roomData"
	"eolian/munc/tcp_server/types"
	"eolian/munc/user"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"log/slog"
	"net"
)

func sendRoomExitMessage(roomName string, conn net.TCPConn, out types.OutChan) {
	exitM := &roomData.RoomExitResp{
		RoomName: roomName,
	}
	exitBytes, err := proto.Marshal(exitM)
	handle_errors.CheckError(err, "Unable to marshal exit announcement data")
	m := &core.PbNetworkResp{
		MuncEventId: constants.TCP_EXIT_ROOM,
		Id:          uuid.NewString(),
		Data: &core.PbNetworkResp_Message{
			Message: exitBytes,
		},
	}
	bytes, err := proto.Marshal(m)
	handle_errors.CheckError(err, "Unable to marshal exit response data")
	out.Channel <- types.OutboundMessage{Bytes: bytes, Conn: conn}
}

func sendTcpMessage(out *types.OutChan) {
	for msg := range out.Channel {
		size := len(msg.Bytes)
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(size))
		_, err := msg.Conn.Write(buf)
		if err != nil {
			slog.Error(err.Error())
		}
		_, err = msg.Conn.Write(msg.Bytes)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func sendObjectSuccess(statusMap map[string]bool, requestId string, conn net.TCPConn, eventId uint32, out types.OutChan) {
	responseData := &roomData.RequestBoolResp{
		Results: statusMap,
	}
	b, err := proto.Marshal(responseData)
	handle_errors.CheckError(err, "Unable to marshal response")
	response := core.PbNetworkResp{
		MuncEventId: eventId,
		Id:          requestId,
		Data: &core.PbNetworkResp_Message{
			Message: b,
		},
	}
	e, err := proto.Marshal(&response)
	handle_errors.CheckError(err, "Unable to marshal response")
	out.Channel <- types.OutboundMessage{Bytes: e, Conn: conn}
}

func broadcastToRoom(sender string, users []user.User, msg []byte, out types.OutChan) {
	for _, usr := range users {
		if usr.UserId == sender {
			continue
		} else {
			out.Channel <- types.OutboundMessage{Bytes: msg, Conn: usr.Socket}
		}
	}
}

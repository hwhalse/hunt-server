package tcp_server

import (
	"eolian/munc/constants"
	handle_errors "eolian/munc/errors"
	"eolian/munc/proto/munc/core"
	"eolian/munc/proto/munc/roomData"
	"eolian/munc/state"
	"eolian/munc/tcp_server/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/anypb"
	"log/slog"
	"net"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func handleTCPInit(message *core.PbNetworkReq, sock net.TCPConn, out *types.OutChan) error {
	log.Info().Msg("TCP Init")
	connInitReq := &roomData.ConnectionInitReq{}
	err := proto.Unmarshal(message.GetMessage(), connInitReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal connection init request")
	}
	if message.UserId == "" {
		userId := roomData.ConnectionInitResp{
			UserId: uuid.New().String(),
		}
		byteArr, _ := proto.Marshal(&userId)
		msg := core.PbNetworkResp{
			MuncEventId: constants.TCP_INIT,
			Message:     byteArr,
		}
		response, err := proto.Marshal(&msg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal response")
		}
		out.Channel <- types.OutboundMessage{Bytes: response, Conn: sock}
		return err
	}
	return err
}

func handleRoomJoin(request *core.PbNetworkReq, conn net.TCPConn, rooms *types.Rooms, out *types.OutChan) error {
	roomJoinData := roomData.RoomConnectReq{}
	err := proto.Unmarshal(request.GetMessage(), &roomJoinData)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal room join request")
	}
	log.Debug().Str("UserId", request.UserId).Str("room", roomJoinData.RoomName).Msg("Received room join request")
	if request.UserId == "" || roomJoinData.RoomName == "" {
		failureMessage := core.MuncFailure{Message: "No userId in room join"}
		encodedFail, err := proto.Marshal(&failureMessage)
		handle_errors.CheckError(err, "Unable to marshal failure msg")
		msg := core.PbNetworkResp{
			MuncEventId: constants.TCP_ERROR,
			Message:     encodedFail,
		}
		send, err := proto.Marshal(&msg)
		handle_errors.CheckError(err, "Unable to marshal fail msg")
		out.Channel <- types.OutboundMessage{Bytes: send, Conn: conn}
		return err
	}

	//add user to room, if no room exists, create one
	rm := rooms.RoomMap[roomJoinData.RoomName]
	if rm == nil {
		err := rooms.AddRoom(roomJoinData.RoomName)
		if err != nil {
			log.Error().Err(err).Msg("Failed to add room")
			return err
		}
		rm = rooms.RoomMap[roomJoinData.RoomName]
		if rm == nil {
			log.Fatal().Msg("Room creation not working")
			return err
		}
	}
	log.Debug().Interface("rm", rm).Msg("Received room join request")
	err = rm.JoinRoom(
		&types.User{UserId: request.UserId, Name: request.UserId},
	)
	if err != nil {
		log.Error().Err(err).Msg("room join error")
		failureMessage := core.MuncFailure{Message: err.Error()}
		encodedFail, err := proto.Marshal(&failureMessage)
		if err != nil {
			log.Error().Err(err).Msg("room failed to marshal failure msg")
		}
		msg := core.PbNetworkResp{
			MuncEventId: constants.TCP_ERROR,
			Message:     encodedFail,
		}
		send, err := proto.Marshal(&msg)
		if err != nil {
			log.Error().Err(err).Msg("room failed to marshal fail msg")
		}
		out.Channel <- types.OutboundMessage{Bytes: send, Conn: conn}
	}
	roomStateData, err := rm.GetData()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get room state data")
		failureMessage := core.MuncFailure{Message: err.Error()}
		encodedFail, err := proto.Marshal(&failureMessage)
		if err != nil {
			log.Error().Err(err).Msg("room failed to marshal failure msg")
		}
		msg := core.PbNetworkResp{
			MuncEventId: constants.TCP_ERROR,
			Message:     encodedFail,
		}
		send, err := proto.Marshal(&msg)
		if err != nil {
			log.Error().Err(err).Msg("room failed to marshal fail msg")
		}
		out.Channel <- types.OutboundMessage{Bytes: send, Conn: conn}
		return err
	}
	var roomStateEvents = make([]*roomData.RoomStateEvt, 0)
	for key, _ := range roomStateData {
		stateObjectArr := make([]*roomData.StateObj, 0)
		x := roomData.RoomStateEvt{
			OwnerId:  "",
			ObjectId: key,
			Payload:  stateObjectArr,
		}
		roomStateEvents = append(roomStateEvents, &x)
	}
	userObjArr := make([]*roomData.StateObj, 0)
	for _, usr := range rm.Users {
		muncUser := &roomData.PbMuncUser{
			UserId: usr.UserId,
			Name:   usr.Name,
		}
		usrAny, err := anypb.New(muncUser)
		if err != nil {
			log.Error().Err(err).Msg("Failed to marshal user object")
			continue
		}
		stateObj := roomData.StateObj{
			Payload: usrAny,
		}
		userObjArr = append(userObjArr, &stateObj)
	}
	userRoomState := roomData.RoomStateEvt{
		Payload: userObjArr,
	}
	roomStateEvents = append(roomStateEvents, &userRoomState)
	roomState := roomData.PbRoomState{
		AllObjects: roomStateEvents,
	}
	responseMsg := roomData.RoomConnectResp{
		RoomName:     roomJoinData.RoomName,
		CurrentState: &roomState,
	}
	responseBytes, _ := proto.Marshal(&responseMsg)
	response := &core.PbNetworkResp{
		MuncEventId: constants.TCP_JOIN_ROOM,
		Message:     responseBytes,
	}
	data, err := proto.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal response msg")
		return err
	}
	out.Channel <- types.OutboundMessage{Bytes: data, Conn: conn}
	return nil
}

func handleRoomExit(request *core.PbNetworkReq, conn net.TCPConn, rooms *types.Rooms, out types.OutChan) {
	usr := state.State.GetUserById(request.UserId)
	msg := roomData.ExitRoomReq{}
	err := proto.Unmarshal(request.Message, &msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal room exit request")
		return
	}
	rm := state.State.GetRoomByName(msg.RoomName)
	slog.Debug("Handle Room Exit", "User", usr.UserId, "Room", msg.RoomName)
	state.State.DeleteFromVoipAddrMap(usr.VoipAddr)
	state.State.DeleteFromPosAddrMap(usr.PositionAddr)
	rm = state.State.GetRoomByName(msg.RoomName)
	ownedIds := rm.GetAllOwnersByUserId(usr.UserId)
	disownResponse := roomData.ReleaseObjectsOwnership{
		ObjectId: ownedIds,
	}
	disownBytes, err := proto.Marshal(&disownResponse)
	handle_errors.CheckError(err, "Unable to marshal disown obj")
	response := core.PbNetworkResp{
		MuncEventId: constants.TCP_DISOWN_OBJECT_BROADCAST,
		Message:     disownBytes,
	}
	responseBytes, _ := proto.Marshal(&response)
	broadcastToRoom(usr.UserId, rm.Users, responseBytes, out)
	rm.DeleteOwnerAll(usr.UserId)
	state.State.RemoveUserFromRoom(msg.RoomName,
		request.UserId)
	sendRoomExitMessage(msg.RoomName, conn, out)
	exitMsg := roomData.UserExitRoomEvt{
		Nickname:     usr.Nickname,
		UserId:       request.UserId,
		RoomName:     rm.Name,
		Disconnected: false,
	}
	exitBytes, err := proto.Marshal(&exitMsg)
	handle_errors.CheckError(err, "Unable to marshal exit announcement data")
	broadcastExit := core.PbNetworkResp{
		MuncEventId: constants.TCP_ANNOUNCE_ROOM_EXIT,
		Message:     exitBytes,
	}
	announceExit, _ := proto.Marshal(&broadcastExit)
	for _, v := range rm.Users {
		out.Channel <- types.OutboundMessage{Bytes: announceExit, Conn: v.Socket}
	}
	state.State.RemoveUser(request.UserId)
}

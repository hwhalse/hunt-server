package main

import (
	"encoding/binary"
	"eolian/munc/constants"
	handle_errors "eolian/munc/errors"
	"eolian/munc/logging"
	"eolian/munc/proto/munc/core"
	"eolian/munc/proto/munc/roomData"
	"eolian/munc/tcp_server"
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var userid = ""
var roomname = "htest"
var testObjId = uuid.NewString()
var nickname = "hh"

var (
	serverAddr     = "localhost:5000"
	testTimeout    = time.Second * 5
	initConnection net.Conn
)

func main() {
	logging.ConfigureLogger()

	// Start the server in a goroutine
	go func() {
		err := tcp_server.StartTCPServer()
		if err != nil {
			log.Fatal().Err(err).Msg("tcp server start error")
		}
	}()
	time.Sleep(3 * time.Second)
	log.Info().Msg("Server started successfully")
	err := initializeServerState()
	if err != nil {
		log.Fatal().Err(err).Msg("initializeServerState error")
	}
	req := roomData.ConnectionInitReq{
		UserId: "",
	}
	reqBytes, err := proto.Marshal(&req)
	init := core.PbNetworkReq{
		MuncEventId: constants.TCP_INIT,
		Id:          uuid.NewString(),
		Data: &core.PbNetworkReq_Message{
			Message: reqBytes,
		},
	}
	bytes, err := proto.Marshal(&init)
	err = sendTcpMessage(bytes)
	if err != nil {
		log.Fatal().Err(err).Msg("sendTcpMessage error")
	}
	for {
		lenBytes := make([]byte, 4)
		_, err := initConnection.Read(lenBytes)
		if err != nil {
			handle_errors.CheckError(err, "Unable to read first bytes")
			if err.Error() == "EOF" {
				fmt.Println("EOF")
				break
			} else {
				return
			}
		}
		length := binary.BigEndian.Uint32(lenBytes)
		protoBytes := make([]byte, length)
		_, err = initConnection.Read(protoBytes)
		muncMessage := &core.PbNetworkResp{}
		err = proto.Unmarshal(protoBytes, muncMessage)
		handle_errors.CheckError(err, "Unable to parse first bytes")
		muncEventId := muncMessage.MuncEventId
		switch muncEventId {
		case constants.TCP_INIT:
			handleTCPInit(muncMessage)
		case constants.TCP_JOIN_ROOM:
			handleRoomJoinMessage(muncMessage)
		case constants.TCP_ROOM_JOIN_SUCCESS:
			handleRoomJoinSuccess(muncMessage)
		case constants.TCP_CREATE_OBJECT:
			handleCreateObjectResponse(muncMessage)
		case constants.TCP_UPDATE_OBJECT:
			handleUpdateObjectResponse(muncMessage)
		case constants.TCP_DISOWN_OBJECT:
			handleDisownObjectResponse(muncMessage)
		case constants.TCP_DELETE_OBJECT:
			handleDeleteObjectResponse(muncMessage)
		case constants.TCP_EXIT_ROOM:
			handleRoomExitResponse()
			break
		}
	}
}

func initializeServerState() error {
	// Establish an initial connection
	conn, err := net.DialTimeout("tcp", serverAddr, testTimeout)
	if err != nil {
		return err
	}

	// Store the connection for potential use in tests
	initConnection = conn

	// Send initialization message
	initMessage := core.PbNetworkReq{
		MuncEventId: constants.TCP_INIT,
		Id:          uuid.New().String(),
	}
	initMessageBytes, err := proto.Marshal(&initMessage)
	err = sendTcpMessage(initMessageBytes)
	if err != nil {
		return err
	}

	return nil
}

func handleRoomExitResponse() {
	log.Info().Msg("Running Room Exit Response")
	err := initConnection.Close()
	handle_errors.CheckError(err, "Unable to close connection")
	log.Info().Msg("Successfully finished all tests")
}

func handleDeleteObjectResponse(response *core.PbNetworkResp) {
	data := &roomData.ObjectUpdateResp{}
	err := proto.Unmarshal(response.GetMessage(), data)
	handle_errors.CheckError(err, "Unable to unmarshal object response")
	log.Info().Bool("success", data.Success).Msg("Object delete response")
	// TODO: Remove this if stattement when UDP testing locally is fixed
	// This will send the Room Exit message and stop the TCP Connection on completion
	sendRoomExitMsg()
}

func sendRoomExitMsg() {
	exitMsg := &roomData.ExitRoomReq{
		RoomName: roomname,
		Nickname: nickname,
	}
	exitMsgBytes, err := proto.Marshal(exitMsg)
	handle_errors.CheckError(err, "Unable to marshal exit room request")
	m := &core.PbNetworkReq{
		MuncEventId: constants.TCP_EXIT_ROOM,
		Id:          testObjId,
		UserId:      userid,
		Data: &core.PbNetworkReq_Message{
			Message: exitMsgBytes,
		},
	}
	bytes, err := proto.Marshal(m)
	handle_errors.CheckError(err, "Unable to marshal exit room request")
	err = sendTcpMessage(bytes)
	if err != nil {
		panic(err)
		return
	}
}

func handleDisownObjectResponse(response *core.PbNetworkResp) {
	data := &roomData.ReleaseObjectsOwnership{}
	err := proto.Unmarshal(response.GetMessage(), data)
	handle_errors.CheckError(err, "Unable to unmarshal object response")
	log.Info().Strs("objectId", data.ObjectId).Msg("object disown response")
	sendObjDeleteMsg()
}

func sendObjDeleteMsg() {
	deleteM := &roomData.ObjectDeleteReq{
		ObjectId: testObjId,
	}
	deleteBytes, err := proto.Marshal(deleteM)
	handle_errors.CheckError(err, "Unable to marshal delete object request")
	m := &core.PbNetworkReq{
		MuncEventId: constants.TCP_DELETE_OBJECT,
		Id:          uuid.NewString(),
		RoomName:    roomname,
		UserId:      userid,
		Data: &core.PbNetworkReq_Message{
			Message: deleteBytes,
		},
	}
	bytes, _ := proto.Marshal(m)
	err = sendTcpMessage(bytes)
	if err != nil {
		panic(err)
		return
	}
}

func handleCreateObjectResponse(response *core.PbNetworkResp) {
	data := &roomData.ObjectUpdateResp{}
	err := proto.Unmarshal(response.GetMessage(), data)
	handle_errors.CheckError(err, "Unable to unmarshal object response")
	log.Info().Bool("success", data.Success).Msg("object create response")
	objectUpdateMsg()
}

func handleUpdateObjectResponse(response *core.PbNetworkResp) {
	data := &roomData.ObjectUpdateResp{}
	err := proto.Unmarshal(response.GetMessage(), data)
	handle_errors.CheckError(err, "Unable to unmarshal object response")
	log.Info().Bool("success", data.Success).Msg("object update response")
	sendDisownMessage()
}

func sendDisownMessage() {
	ids := []string{testObjId}
	disown := &roomData.ReleaseObjectsOwnership{
		ObjectId: ids,
	}
	disownBytes, err := proto.Marshal(disown)
	handle_errors.CheckError(err, "Unable to marshal disown message")
	m := &core.PbNetworkReq{
		MuncEventId: constants.TCP_DISOWN_OBJECT,
		Id:          uuid.NewString(),
		RoomName:    roomname,
		UserId:      userid,
		Data: &core.PbNetworkReq_Message{
			Message: disownBytes,
		},
	}
	bytes, err := proto.Marshal(m)
	handle_errors.CheckError(err, "Unable to marshal disown message")
	err = sendTcpMessage(bytes)
	if err != nil {
		panic(err)
		return
	}
}

func objectUpdateMsg() {
	name := "updated"
	avatar := roomData.PbMuncAvatar{
		UserId: userid,
		Name:   &name,
	}
	avatarBytes, _ := anypb.New(&avatar)
	payloadObj := &roomData.StateObj{
		Payload: avatarBytes,
	}
	var payloadArr = make([]*roomData.StateObj, 1)
	payloadArr = append(payloadArr, payloadObj)
	objEvent := &roomData.RoomStateEvt{
		ObjectId: testObjId,
		Payload:  payloadArr,
	}
	objBytes, err := proto.Marshal(objEvent)
	handle_errors.CheckError(err, "Unable to marshal objEvent")
	m := &core.PbNetworkReq{
		MuncEventId: constants.TCP_UPDATE_OBJECT,
		Id:          uuid.NewString(),
		RoomName:    roomname,
		UserId:      userid,
		Data: &core.PbNetworkReq_Message{
			Message: objBytes,
		},
	}
	fullBytes, _ := proto.Marshal(m)
	err = sendTcpMessage(fullBytes)
	if err != nil {
		panic(err)
		return
	}
}

func handleRoomJoinMessage(muncMessage *core.PbNetworkResp) {
	data := &roomData.RoomConnectResp{}
	err := proto.Unmarshal(muncMessage.GetMessage(), data)
	handle_errors.CheckError(err, "Unable to parse room connect")
	log.Info().Str("state", data.CurrentState.String()).Msg("Room State on join")
	sendRoomJoinReq()
}

func handleTCPInit(response *core.PbNetworkResp) {
	log.Info().Msg(response.String())
	initResponse := roomData.ConnectionInitResp{}
	err := proto.Unmarshal(response.GetMessage(), &initResponse)
	handle_errors.CheckError(err, "Unable to get message from init response")
	userid = initResponse.UserId
	sendRoomConnectReq()
}

func handleRoomJoinSuccess(response *core.PbNetworkResp) {
	log.Info().Msg(response.String())
	data := roomData.RoomJoinResp{}
	err := proto.Unmarshal(response.GetMessage(), &data)
	handle_errors.CheckError(err, "Unable to get message from room connect response")
	// TODO: Fix UDP testing locally

	sendCreateObject()
}

func sendCreateObject() {
	name := "test"
	platform := roomData.Platform_DESKTOP
	var payloadObjects = make([]*roomData.StateObj, 1)
	avatar := &roomData.PbMuncAvatar{
		UserId:   userid,
		Name:     &name,
		Platform: &platform,
	}
	avatarBytes, _ := anypb.New(avatar)
	p := &roomData.StateObj{
		Payload: avatarBytes,
	}
	payloadObjects = append(payloadObjects, p)
	object := roomData.RoomStateEvt{
		ObjectId: testObjId,
		Payload:  payloadObjects,
	}
	objBytes, err := proto.Marshal(&object)
	handle_errors.CheckError(err, "Unable to marshal object")
	msg := &core.PbNetworkReq{
		MuncEventId: constants.TCP_CREATE_OBJECT,
		RoomName:    roomname,
		UserId:      userid,
		Id:          uuid.NewString(),
		Data: &core.PbNetworkReq_Message{
			Message: objBytes,
		},
	}
	bytes, err := proto.Marshal(msg)
	handle_errors.CheckError(err, "Unable to marshal object")
	err = sendTcpMessage(bytes)
	if err != nil {
		panic(err)
		return
	}
}

func sendRoomJoinReq() {
	reqData := roomData.RoomJoinReq{
		VoipPort:     5011,
		PositionPort: 5010,
		Platform:     roomData.Platform_HOLOLENS,
	}
	dB, err := proto.Marshal(&reqData)
	handle_errors.CheckError(err, "unable to marshal")
	reqM := core.PbNetworkReq{
		MuncEventId: constants.TCP_ROOM_JOIN_SUCCESS,
		Id:          uuid.NewString(),
		RoomName:    roomname,
		UserId:      userid,
		Data: &core.PbNetworkReq_Message{
			Message: dB,
		},
	}
	b, err := proto.Marshal(&reqM)
	handle_errors.CheckError(err, "unable to marshal")
	err = sendTcpMessage(b)
	if err != nil {
		panic(err)
		return
	}
}

func sendRoomConnectReq() {
	roomConnect := roomData.RoomConnectReq{
		RoomName: roomname,
		Nickname: nickname,
	}
	conBytes, err := proto.Marshal(&roomConnect)
	handle_errors.CheckError(err, "Unable to marshal roomConnect data")
	netReq := core.PbNetworkReq{
		MuncEventId: constants.TCP_JOIN_ROOM,
		UserId:      userid,
		Id:          uuid.NewString(),
		Data:        &core.PbNetworkReq_Message{Message: conBytes},
	}
	reqBytes, _ := proto.Marshal(&netReq)
	err = sendTcpMessage(reqBytes)
	if err != nil {
		panic(err)
		return
	}
}

func sendTcpMessage(msg []byte) error {
	size := len(msg)
	log.Debug().Int("size", size).Msg("Length of bytes")
	buf := make([]byte, 4)
	log.Debug().Bytes("buf", buf).Msg("Buffer Before")
	binary.BigEndian.PutUint32(buf, uint32(size))
	log.Debug().Bytes("buf", buf).Msg("Buffer After")
	_, err := initConnection.Write(buf)
	if err != nil {
		return err
	}
	log.Info().Msg("writing msg")
	_, err = initConnection.Write(msg)
	if err != nil {
		return err
	}
	return nil
}

package main

import (
	"encoding/binary"
	"eolian/munc/constants"
	handle_errors "eolian/munc/errors"
	"eolian/munc/proto/munc/core"
	"eolian/munc/proto/munc/roomData"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var userid = ""
var conn *net.TCPConn
var roomname = "htest"
var testObjId = uuid.NewString()
var nickname = "hh"

func main() {
	cloudAddr := "esp-dev.eastus2.cloudapp.azure.com:5000"
	cloudAddress, err := net.ResolveTCPAddr("tcp", cloudAddr)
	handle_errors.CheckError(err, "Unable to resolve cloud tcp address")
	tcpConn, err := net.DialTCP("tcp", nil, cloudAddress)
	conn = tcpConn
	handle_errors.CheckError(err, "Unable to connect to server")
	defer tcpConn.Close()
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
	sendTcpMessage(bytes)
	for {
		lenBytes := make([]byte, 4)
		_, err := tcpConn.Read(lenBytes)
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
		_, err = tcpConn.Read(protoBytes)
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

func handleRoomExitResponse() {
	fmt.Println("Room Exit Response")
	//err := conn.Close()
	//handle_errors.CheckError(err, "Unable to close connection")
}

func handleDeleteObjectResponse(response *core.PbNetworkResp) {
	data := &roomData.ObjectUpdateResp{}
	err := proto.Unmarshal(response.GetMessage(), data)
	handle_errors.CheckError(err, "Unable to unmarshal object response")
	fmt.Println("object delete response", data.Success)
	//sendRoomExitMsg()
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
	sendTcpMessage(bytes)
}

func handleDisownObjectResponse(response *core.PbNetworkResp) {
	data := &roomData.ReleaseObjectsOwnership{}
	err := proto.Unmarshal(response.GetMessage(), data)
	handle_errors.CheckError(err, "Unable to unmarshal object response")
	fmt.Println("object disown response", data.ObjectId)
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
	bytes, err := proto.Marshal(m)
	sendTcpMessage(bytes)
}

func handleCreateObjectResponse(response *core.PbNetworkResp) {
	data := &roomData.ObjectUpdateResp{}
	err := proto.Unmarshal(response.GetMessage(), data)
	handle_errors.CheckError(err, "Unable to unmarshal object response")
	fmt.Println("object create response", data.Success)
	objectUpdateMsg()
}

func handleUpdateObjectResponse(response *core.PbNetworkResp) {
	data := &roomData.ObjectUpdateResp{}
	err := proto.Unmarshal(response.GetMessage(), data)
	handle_errors.CheckError(err, "Unable to unmarshal object response")
	fmt.Println("object update response", data.Success)
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
	sendTcpMessage(bytes)
}

func objectUpdateMsg() {
	platform := roomData.Platform_HOLOLENS
	avatar := &roomData.PbMuncAvatar{
		UserId:   "123",
		Platform: &platform,
	}
	avatarBytes, err := anypb.New(avatar)
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
	fullBytes, err := proto.Marshal(m)
	sendTcpMessage(fullBytes)
}

func handleRoomJoinMessage(muncMessage *core.PbNetworkResp) {
	data := &roomData.RoomConnectResp{}
	err := proto.Unmarshal(muncMessage.GetMessage(), data)
	handle_errors.CheckError(err, "Unable to parse room connect")
	fmt.Println("state", data.CurrentState)
	sendRoomJoinReq()
}

func handleTCPInit(response *core.PbNetworkResp) {
	fmt.Println(response)
	initResponse := roomData.ConnectionInitResp{}
	err := proto.Unmarshal(response.GetMessage(), &initResponse)
	handle_errors.CheckError(err, "Unable to get message from init response")
	userid = initResponse.UserId
	sendRoomConnectReq()
}

func handleRoomJoinSuccess(response *core.PbNetworkResp) {
	fmt.Println(response)
	data := roomData.RoomJoinResp{}
	err := proto.Unmarshal(response.GetMessage(), &data)
	handle_errors.CheckError(err, "Unable to get message from room connect response")
	fmt.Println("udp ports", data.PositionServerPort, data.VoipServerPort)
	go createUdpServers(data.PositionServerPort, data.VoipServerPort)
	sendCreateObject()
}

func createUdpServers(posPort int32, voipPort int32) {
	slog.Debug("Create UDP", "Server", "one")
	udpPos, err := net.ListenPacket("udp4", "0.0.0.0:6010")
	handle_errors.CheckError(err, "Unable to start pos server")
	go readUdp(udpPos)
	handle_errors.CheckError(err, "Unable to send udp")
	udpVoip, err := net.ListenPacket("udp4", "0.0.0.0:6011")
	handle_errors.CheckError(err, "Unable to start voip server")
	slog.Debug("UDP server", "status", udpPos.LocalAddr())
	go readUdp(udpVoip)
	go sendUdpMessages(udpPos, udpVoip)
}

func readUdp(c net.PacketConn) {
	buf := make([]byte, 1024)
	for {
		n, addr, err := c.ReadFrom(buf)
		t := time.Now()
		err = t.GobDecode(buf[:n])
		c := time.Since(t)
		fmt.Println(c, addr)
		handle_errors.CheckError(err, "Unable to read from connection")
	}
}

func sendUdpMessages(pos net.PacketConn, voip net.PacketConn) {
	a, err := net.ResolveUDPAddr("udp4", "esp-dev.eastus2.cloudapp.azure.com:43234")
	v, err := net.ResolveUDPAddr("udp4", "esp-dev.eastus2.cloudapp.azure.com:44234")
	t := time.Now()
	b, err := t.GobEncode()
	_, err = pos.WriteTo(b, a)
	handle_errors.CheckError(err, "Unable to send position")
	_, err = voip.WriteTo(b, v)
	handle_errors.CheckError(err, "Unable to send voice")
	time.Sleep(10 * time.Second)
	sendUdpMessages(pos, voip)
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
	avatarBytes, err := anypb.New(avatar)
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
	sendTcpMessage(bytes)
}

func sendRoomJoinReq() {
	reqData := roomData.RoomJoinReq{
		VoipPort:     6011,
		PositionPort: 6010,
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
	sendTcpMessage(b)
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
	reqBytes, err := proto.Marshal(&netReq)
	sendTcpMessage(reqBytes)
}

func sendTcpMessage(msg []byte) {
	size := len(msg)
	fmt.Println("size ", size)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(size))
	_, err := conn.Write(buf)
	handle_errors.CheckError(err, "Unable to write binary delim")
	fmt.Println("writing msg")
	_, err = conn.Write(msg)
	handle_errors.CheckError(err, "Unable to write tcp join room response")
}

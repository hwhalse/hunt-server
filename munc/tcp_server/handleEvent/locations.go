package handleEvent

import (
	"eolian/munc/constants"
	"eolian/munc/proto/munc/core"
	"eolian/munc/tcp_server/types"
	"github.com/golang/protobuf/proto"
	"net"
)

func GetLocations(request *core.PbNetworkReq, conn net.TCPConn, rooms *types.Rooms, out *types.OutChan) error {
	room := rooms.RoomMap[request.RoomName]
	locations, err := room.GetUserLocations()
	if err != nil {
		return err
	}
	locationProtos := make([]*core.Location, len(*locations))
	for i, location := range *locations {
		locationProtos[i] = &core.Location{
			Lat: location.Lat,
			Lon: location.Lon,
			Alt: location.Alt,
		}
	}
	locationMsg := core.Locations{
		Locations: locationProtos,
	}
	locationMsgBytes, err := proto.Marshal(&locationMsg)
	response := core.PbNetworkResp{
		MuncEventId: constants.GET_LOCATIONS,
		Message:     locationMsgBytes,
	}
	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		return err
	}
	out.Channel <- types.OutboundMessage{Bytes: responseBytes, Conn: conn}
	return nil
}

func UpdateLocation(request *core.PbNetworkReq, conn net.TCPConn, rooms *types.Rooms, out *types.OutChan) error {
	room := rooms.RoomMap[request.RoomName]
	user, err := room.GetUserById(request.UserId)
	if err != nil {
		return err
	}
	locationFromProto := &core.Location{}
	err = proto.Unmarshal(request.Message, locationFromProto)
	if err != nil {
		return err
	}
	user.Location = types.Location{
		Lat: locationFromProto.Lat,
		Lon: locationFromProto.Lon,
		Alt: locationFromProto.Alt,
	}
	return nil
}

package state

import (
    "eolian/munc/proto/munc/roomData"
	"eolian/munc/room"
	"eolian/munc/user"
	"net"
)

var State GameState
var DeletedUsers DeleteQueue

type MuncObject struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

type DeleteQueue struct {
	Queue []user.User
}

func init() {
	State.PositionServerRooms = map[string]string{}
	State.VoipRooms = map[string]string{}
}

func (q *DeleteQueue) AddToQueue(usr user.User) {
	q.Queue = append(q.Queue, usr)
}

func (q *DeleteQueue) RemoveFromQueue(userId string) {
	var replacement []user.User
	for _, v := range q.Queue {
		if v.UserId != userId {
			replacement = append(replacement, v)
		}
	}
	q.Queue = replacement
}

func (q *DeleteQueue) RemoveFromQueueByAddr(addr net.TCPConn) {
	var replacement []user.User
	for _, v := range q.Queue {
		if v.Socket.RemoteAddr() != addr.RemoteAddr() {
			replacement = append(replacement, v)
		}
	}
	q.Queue = replacement
}

func (q *DeleteQueue) UserInQueue(userId string) bool {
	for _, v := range q.Queue {
		if v.UserId == userId {
			return true
		}
	}
	return false
}

func (q *DeleteQueue) CheckAddrInDeleteQueue(addr net.TCPConn, userId string) {
	for _, v := range q.Queue {
		if v.Socket.RemoteAddr() == addr.RemoteAddr() {
			q.RemoveFromQueueByAddr(addr)
			State.OverwriteId(userId, addr)
		}
	}
}

func (state *GameState) OverwriteId(userId string, addr net.TCPConn) {
	for _, v := range state.Users {
		if v.Socket.RemoteAddr() == addr.RemoteAddr() {
			v.UserId = userId
		}
	}
}

type GameState struct {
	Objects             []MuncObject
	Rooms               []room.Room
	Users               []user.User
	TCPConns            []net.TCPConn
	PositionServerRooms map[string]string
	VoipRooms           map[string]string
}

func (state *GameState) GetUserByConn(conn net.TCPConn) user.User {
	for _, v := range state.Users {
		if v.Socket == conn {
			return v
		}
	}
	return user.User{}
}

func (state *GameState) UpdateUser(usr user.User) {
	for i, v := range state.Users {
		if v.UserId == usr.UserId {
			state.Users[i] = usr
		}
	}
}

func (state *GameState) UpdateUserPorts(userId string, positionAddr string, voiceAddr string) {
	for i, v := range state.Users {
		if v.UserId == userId {
			state.Users[i].PositionAddr = positionAddr
			state.Users[i].VoipAddr = voiceAddr
		}
	}
}

func (state *GameState) UpdateUserPlatform(userId string, platform roomData.Platform) {
	for i, v := range state.Users {
		if v.UserId == userId {
			state.Users[i].Platform = platform
		}
	}
}

func (state *GameState) UpdateUserConnected(userId string, status bool) {
	for i, v := range state.Users {
		if v.UserId == userId {
			state.Users[i].Connected = status
		}
	}
}

func (state *GameState) RemoveUser(userId string) {
	var replacement []user.User
	for _, v := range state.Users {
		if v.UserId != userId {
			replacement = append(replacement, v)
		}
	}
	state.Users = replacement
}

func (state *GameState) GetUserById(userId string) user.User {
	for _, v := range state.Users {
		if v.UserId == userId {
			return v
		}
	}
	return user.User{}
}

func (state *GameState) AddTCPConn(conn net.TCPConn) {
	conns := make([]net.TCPConn, len(state.TCPConns))
	copy(conns, state.TCPConns)
	conns = append(conns, conn)
	state.TCPConns = conns
	return
}

func (state *GameState) UpdateUserByIndex(update user.User, index int) {
	state.Users[index] = update
}

func (state *GameState) FindUserById(id string) (user.User, int) {
	for i, v := range state.Users {
		if v.UserId == id {
			return v, i
		}
	}
	return user.User{}, -1
}

func (state *GameState) AddUser(usr user.User) {
	if len(state.Users) == 0 {
		state.Users = append(state.Users, usr)
	} else {
		usrs := make([]user.User, len(state.Users))
		copy(usrs, state.Users)
		usrs = append(usrs, usr)
		state.Users = usrs
	}
}

func (state *GameState) AddRoom(newRoom room.Room) {
	rooms := make([]room.Room, len(state.Rooms))
	copy(rooms, state.Rooms)
	rooms = append(rooms, newRoom)
	state.Rooms = rooms
	return
}

func (state *GameState) AddUserToRoom(user user.User, roomName string, password string) bool {
	for i, v := range state.Rooms {
		if v.Name == roomName {
			if state.Rooms[i].Password != "" {
				if state.Rooms[i].Password != password {
					return false
				}
			}
			v.Users = append(v.Users, user)
			state.Rooms[i] = v
			return true
		}
	}
	return false
}

func (state *GameState) GetRoomByName(name string) room.Room {
	for _, v := range state.Rooms {
		if v.Name == name {
			return v
		}
	}
	return room.Room{}
}

func (state *GameState) GetRoomByPosAddr(addr string) room.Room {
	roomName := state.PositionServerRooms[addr]
	return state.GetRoomByName(roomName)
}

func (state *GameState) DeleteFromPosAddrMap(addr string) {
	delete(state.PositionServerRooms, addr)
}

func (state *GameState) DeleteFromVoipAddrMap(addr string) {
	delete(state.VoipRooms, addr)
}

func (state *GameState) GetRoomByVoipAddr(addr string) room.Room {
	roomName := state.VoipRooms[addr]
	return state.GetRoomByName(roomName)
}

func (state *GameState) AddPosAddrToRoom(addr string, roomName string) {
	_, ok := state.PositionServerRooms[addr]
	if ok {
		return
	}
	state.PositionServerRooms[addr] = roomName
}

func (state *GameState) AddVoipAddrToRoom(addr string, roomName string) {
	state.VoipRooms[addr] = roomName
}

func (state *GameState) RemoveUserFromState(userId string) {
	var replacement []user.User
	for _, v := range state.Users {
		if v.UserId != userId {
			replacement = append(replacement, v)
		}
	}
	state.Users = replacement
}

func (state *GameState) RemoveUserFromRoom(roomName string, userId string) {
	for i, v := range state.Rooms {
		if v.Name == roomName {
			v.Users = v.RemoveUserFromRoom(userId)
			state.Rooms[i] = v
		}
	}
}

func (state *GameState) GetObjs() []MuncObject {
	return state.Objects
}

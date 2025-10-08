package room

import (
	"eolian/munc/proto/munc/roomData"
	"eolian/munc/user"
	"net"
	"sync"
)

type Room struct {
	Name       string
	Users      []user.User
	State      any
	Objects    *SafeMap
	Owners     *SafeStringMap
	UpdateLock *sync.Mutex
	CreateLock *sync.Mutex
	DeleteLock *sync.Mutex
	Password   string
}

type SafeMap struct {
	sync.RWMutex
	writeChan <-chan WriteObject
	data      map[string]*roomData.StateObj
}

type SafeStringMap struct {
	sync.RWMutex
	data map[string]string
}

type MuncDataObject struct {
	lock    sync.RWMutex
	payload any
}

func NewRoom(name string) *Room {
	return &Room{
		Name:       name,
		Objects:    NewSafeMap(),
		Owners:     NewSafeStringMap(),
		UpdateLock: &sync.Mutex{},
		CreateLock: &sync.Mutex{},
		DeleteLock: &sync.Mutex{},
	}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		data:      make(map[string]*roomData.StateObj),
		writeChan: make(<-chan WriteObject),
	}
}

func NewSafeStringMap() *SafeStringMap {
	return &SafeStringMap{data: make(map[string]string)}
}

type WriteObject struct {
	ObjectId string
	OwnerId  string
	Data     any
}

func (s *SafeMap) Get(key string) (any, bool) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *SafeMap) GetAllData() map[string]*roomData.StateObj {
	s.RLock()
	defer s.RUnlock()
	return s.data
}

func (s *SafeMap) Set(key string, val *roomData.StateObj) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = val
}

func (s *SafeMap) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, key)
}

func (s *SafeMap) Contains(key string) bool {
	s.Lock()
	defer s.Unlock()
	val, ok := s.data[key]
	return ok && val != nil
}

func (s *SafeStringMap) Get(key string) (any, bool) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *SafeStringMap) GetAllData() map[string]string {
	s.RLock()
	defer s.RUnlock()
	return s.data
}

func (s *SafeStringMap) Set(key string, val string) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = val
}

func (s *SafeStringMap) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, key)
}

func (s *SafeStringMap) Contains(key string) bool {
	s.Lock()
	defer s.Unlock()
	val, ok := s.data[key]
	return ok && val != ""
}

func (room *Room) AddUserToRoom(newUser user.User) {
	room.Users = append(room.Users, newUser)
	return
}

func (room *Room) GetData() map[string]*roomData.StateObj {
	return room.Objects.GetAllData()
}

func (room *Room) AddObject(objectId string, payload *roomData.StateObj, ownerId string) bool {
	if room.Owners.Contains(objectId) {
		return false
	}
	if ownerId != "" {
		room.Owners.Set(objectId, ownerId)
	}
	room.Objects.Set(objectId, payload)
	return true
}

func (room *Room) DeleteOwner(objectId string, requestUserId string) bool {
	owner, _ := room.Owners.Get(objectId)
	if requestUserId != owner {
		return false
	} else {
		room.Owners.Delete(objectId)
		return true
	}
}

func (room *Room) DeleteOwnerAll(requestUserId string) {
	for objectId, owner := range room.Owners.data {
		if requestUserId == owner {
			room.Owners.Set(objectId, "")
		}
	}
}

func (room *Room) GetAllOwnersByUserId(userId string) []string {
	objectIdArray := make([]string, 0)
	for objectId, owner := range room.Owners.data {
		if userId == owner {
			objectIdArray = append(objectIdArray, objectId)
		}
	}
	return objectIdArray
}

func (room *Room) DeleteObject(objectId string, requestUserId string) bool {
	owner, _ := room.Owners.Get(objectId)
	if requestUserId == owner {
		room.Objects.Delete(objectId)
		room.Owners.Delete(objectId)
		return true
	} else if owner == nil {
		room.Objects.Delete(objectId)
		return true
	} else {
		return false
	}
}

func (room *Room) RemoveObject(objectId string) {
	room.Objects.Delete(objectId)
	room.Owners.Delete(objectId)
}

func (room *Room) UpdateOwner(objectId string, ownerId string) {
	room.Owners.Set(objectId, ownerId)
}

func (room *Room) GetOwner(objectId string) string {
	val, retrieved := room.Owners.Get(objectId)
	if retrieved {
		return val.(string)
	} else {
		return ""
	}
}

func (room *Room) RemoveOwner(objectId string) {
	room.Owners.Set(objectId, "")
}

func (room *Room) UpdateObject(objectId string, payload *roomData.StateObj, userId string, ownerId string) bool {
	room.UpdateLock.Lock()
	defer room.UpdateLock.Unlock()
	owner := room.GetOwner(objectId)
	if ownerId != "" {
		if owner == ownerId {
			room.Objects.Set(objectId, payload)
			return true
		} else if owner != "" && owner != userId {
			return false
		} else {
			room.Objects.Set(objectId, payload)
			room.UpdateOwner(objectId, ownerId)
			return true
		}
	} else {
		if owner != "" {
			return false
		} else {
			room.Objects.Set(objectId, payload)
			return true
		}
	}
}

func (room *Room) OverwriteId(userId string, sock net.Conn) bool {
	for _, v := range room.Users {
		if v.Socket.RemoteAddr().String() == sock.RemoteAddr().String() {
			v.UserId = userId
			return true
		}
	}
	return false
}

func (room *Room) RemoveUserFromRoom(userId string) []user.User {
	var replacement []user.User
	for _, v := range room.Users {
		if v.UserId != userId {
			replacement = append(replacement, v)
		}
	}
	return replacement
}

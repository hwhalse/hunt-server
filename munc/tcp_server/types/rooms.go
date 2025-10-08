package types

import (
	"errors"
	"sync"
)

type Rooms struct {
	sync.Mutex
	RoomMap map[string]*Room
}

func (rooms *Rooms) AddRoom(name string) error {
	if _, ok := rooms.RoomMap[name]; ok {
		return errors.New("room already exists")
	}
	rooms.Lock()
	defer rooms.Unlock()
	rooms.RoomMap[name] = &Room{}
	return nil
}

func (rooms *Rooms) removeRoom(name string) error {
	if rooms.RoomMap[name] == nil {
		return errors.New("room Not Found")
	}
	delete(rooms.RoomMap, name)
	return nil
}

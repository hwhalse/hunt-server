package state

import (
	"errors"
	"hunt/models"
	"sync"
)

const maxClients = 256

type HuntState struct {
	mu sync.RWMutex
	Uids [maxClients]string
	Callsigns [maxClients]string
	Free [maxClients]bool
}

type Locations struct {
	mu sync.RWMutex
	Lat [maxClients]float64
	Lon [maxClients]float64
	Alt [maxClients]float64
	Uid [maxClients]string
	Free [maxClients]bool
}

func NewStateObject() *HuntState {
	return &HuntState{
		Uids: [maxClients]string{},
		Callsigns: [maxClients]string{},
		Free: [maxClients]bool{},
	}
}

func NewLocationsObject() *Locations {
	return &Locations{
		Lat:  [maxClients]float64{},
		Lon:  [maxClients]float64{},
		Alt:  [maxClients]float64{},
		Uid:  [maxClients]string{},
		Free: [maxClients]bool{},
	}
}

func (l *Locations) GetUserIndex(uid string) int {
	for i := range l.Uid {
		if l.Uid[i] == uid {
			return i
		}
	}
	return -1
}

func (l *Locations) AddLocation(lat, lon, alt float64, uid string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	freeIndex := l.getFreeIndex()
	if freeIndex < 0 {
		return errors.New("Locations is full")
	}
	l.Uid[freeIndex] = uid
	l.Free[freeIndex] = true
	l.Lat[freeIndex] = lat
	l.Lon[freeIndex] = lon
	l.Alt[freeIndex] = alt
	return nil
}

func (l *Locations) getFreeIndex() int {
	for i := range l.Free {
		if l.Free[i] == false {
			return i
		}
	}
	return -1
}

func (l *Locations) UpdateLocation(lat, lon, alt float64, uid string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	index := l.GetUserIndex(uid)
	if index < 0 {
		return errors.New("Unable to find user")
	}
	l.Lat[index] = lat
	l.Lon[index] = lon
	l.Alt[index] = alt
	return nil
}

func (l *Locations) GetLocation(index int) *models.Location {
	return &models.Location{
		Lat: l.Lat[index],
		Lon: l.Lon[index],
		Alt: l.Alt[index],
	}
}

func (l *Locations) GetLocationByUid(uid string) (*models.Location, error) {
	index := -1
	for i := range l.Uid {
		if l.Uid[i] == uid {
			index = i
		}
	}
	if index < 0 {
		return nil, errors.New("Unable to find user in locations")
	}
	return &models.Location{
		Lat: l.Lat[index],
		Lon: l.Lon[index],
		Alt: l.Alt[index],
	}, nil
}

func (l *Locations) DeleteLocation(index int) {
	l.Lat[index] = 0.0
	l.Lon[index] = 0.0
	l.Alt[index] = 0.0
	l.Free[index] = false
}

func (s *HuntState) AddClient(uid, callsign string, location models.Location) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	freeSpace := s.getFreeIndex()
	if freeSpace < 0 {
		return -1, errors.New("Server capacity reached")
	}
	s.Free[freeSpace] = true
	s.Callsigns[freeSpace] = callsign
	s.Uids[freeSpace] = uid
	return freeSpace, nil
}

func (s *HuntState) RemoveClient(index int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Free[index] = false
	s.Uids[index] = ""
	s.Callsigns[index] = ""
}

func (s *HuntState) FindClient(uid string) (models.HuntUser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.Uids {
		if s.Uids[i] == uid {
			return models.HuntUser{
				Callsign: s.Callsigns[i],
				Uid: s.Uids[i],
			}, nil
		}
	}
	return models.HuntUser{}, errors.New("Unable to find client")
}

func (s *HuntState) FindClientIndex(uid string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.Uids {
		if s.Uids[i] == uid {
			return i
		}
	}
	return -1
}

func (s *HuntState) getFreeIndex() int {
	for i := range s.Free {
		if !s.Free[i] {
			return i
		}
	}
	return -1
}
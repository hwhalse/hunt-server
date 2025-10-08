package types

import (
	"errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"sync"
)

type Room struct {
	sync.RWMutex
	Name  string
	Users map[string]*User
	State map[string]*Target
}

func (r *Room) JoinRoom(user *User) error {
	log.Debug().Interface("user", *user).Msg("join room")
	r.Lock()
	defer r.Unlock()
	if _, ok := r.Users[user.Name]; ok {
		return errors.New("user exists")
	}
	r.Users[user.UserId] = &User{UserId: user.UserId, Name: user.Name}
	return nil
}

func (r *Room) LeaveRoom(user *User) error {
	log.Debug().Interface("user", *user).Msg("leave room")
	r.Lock()
	defer r.Unlock()
	if _, ok := r.Users[user.UserId]; ok {
		delete(r.Users, user.UserId)
		return nil
	}
	return errors.New("user does not exist")
}

func (r *Room) GetData() (map[string]*Target, error) {
	if len(r.Users) == 0 {
		return nil, errors.New("no players")
	}
	return r.State, nil
}

func (r *Room) AddTarget(target *Target) {
	r.Lock()
	defer r.Unlock()
	if target.Uid == "" {
		target.Uid = uuid.New().String()
	}
	if r.State == nil {
		r.State = make(map[string]*Target)
	}
	r.State[target.Uid] = target
}

func (r *Room) RemoveTarget(target *Target) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.State[target.Uid]; ok {
		delete(r.State, target.Uid)
		return nil
	}
	return errors.New("target does not exist")
}

func (r *Room) AddUser(user *User) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.Users[user.UserId]; ok {
		return errors.New("user exists")
	}
	r.Users[user.UserId] = user
	return nil
}

func (r *Room) RemoveUser(user *User) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.Users[user.UserId]; ok {
		delete(r.Users, user.UserId)
		return nil
	}
	return errors.New("user does not exist")
}

func (r *Room) GetUserLocations() (*[]Location, error) {
	r.RLock()
	defer r.RUnlock()
	locations := make([]Location, 0, len(r.Users))
	if len(r.Users) == 0 {
		return nil, errors.New("no players")
	}
	for _, user := range r.Users {
		location, err := user.GetLocation()
		if err != nil {
			log.Error().Interface("user", user).Err(err).Msg("invalid user location")
			continue
		}
		locations = append(locations, location)
	}
	return &locations, nil
}

func (r *Room) GetUserById(userId string) (*User, error) {
	r.RLock()
	defer r.RUnlock()
	if u, ok := r.Users[userId]; ok {
		return u, nil
	}
	return nil, errors.New("user does not exist")
}

package types

import "errors"

type User struct {
	UserId   string
	Name     string
	Location Location
}

func (u *User) GetLocation() (Location, error) {
	if u.Location.Lon == 0.0 && u.Location.Lat == 0.0 && u.Location.Alt == 0.0 {
		return Location{}, errors.New("location is empty")
	}
	return u.Location, nil
}

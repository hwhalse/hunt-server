package socket

import (
	"context"
	"encoding/json"
	"hunt/collections"
	"hunt/models"
	"log"

	"github.com/google/uuid"
	"github.com/lxzan/gws"
)

func handleInit(e *Event, socket *gws.Conn, useState bool) {
	payload := models.InitMsg{}
	err := json.Unmarshal([]byte(e.Payload), &payload)
	if err != nil {
		log.Fatal(err)
	}
	newUid := uuid.NewString()

	responseEvent := &Event{
		Type:    NewUid,
		Payload: newUid, 
		Time: e.Time,
	}
	b, err := json.Marshal(responseEvent)
	if err != nil {
		log.Fatalln("unable to marshal response", err)
		socket.NetConn().Close()
		return
	}
	if payload.Uid == "" {
		newUser := models.HuntUser{
			Callsign: payload.Callsign,
			Uid: newUid,
			Location: models.Location{},
			Active: true,
		}
		connMap.Store(socket, newUser)
		if useState {
			addNewUserToState(newUser)
		} else {
			addNewUserToDb(newUser)
		}
		socket.WriteAsync(gws.OpcodeText, b, func(err error) {
			if err != nil {
				socket.NetConn().Close()
			}
		})
	}
	
	currentState := models.InitResponse{}

	var users []models.HuntUser
	if useState {
		users = getUsersFromState()
	} else {
		users = getUsersFromDB()
	}
	currentState.Users = users
	stateToBytes, err := json.Marshal(currentState)
	if err != nil {
		log.Fatalf("Unable to marshal state init %e\n", err)
	}
	stateEvent := &Event{
		Type: InitResponse,
		Payload: string(stateToBytes),
		Time: e.Time,
	}
	eventBytes, err := json.Marshal(stateEvent)
	if err != nil {
		log.Fatalf("Unable to marshal state init %e\n", err)
	}
	socket.WriteAsync(gws.OpcodeText, eventBytes, func(err error) {
		if err != nil {
			log.Fatalf("Unable to send state init %e\n", err)
		}
	})
}

func addNewUserToState(newUser models.HuntUser) {	
	_, err := huntState.AddClient(newUser.Uid, newUser.Callsign, newUser.Location)
	if err != nil {
		log.Fatalln("Unable to add user to state", err)
	}
	err = locations.AddLocation(0.0, 0.0, 0.0, newUser.Uid)
	if err != nil {
		log.Fatalln("Unable to add user location", err)
	}
}

func addNewUserToDb(newUser models.HuntUser) {
	err := collections.UsersCollectionManager.AddUser(context.Background(), newUser)
	if err != nil {
		log.Fatalf("Unable to add user during init %e\n", err)
	}
}

func getUsersFromState() []models.HuntUser {
	allUsers := []models.HuntUser{}
	for i := range huntState.Uids {
		if huntState.Free[i] {
			uid := huntState.Uids[i]
			huntUser, err := huntState.FindClient(uid)
			if err != nil {
				log.Fatalln("Unable to find hunt user", err)
			}
			loc, err := locations.GetLocationByUid(uid)
			if err != nil {
				log.Fatalln("User location does not exist", err)
			}
			huntUser.Location = *loc
			allUsers = append(allUsers, huntUser)
		}
	}
	return allUsers
}

func getUsersFromDB() []models.HuntUser {
	users, err := collections.UsersCollectionManager.FindAll(context.Background())
	if err != nil {
		log.Fatalln("Unable to get locations", err)
	}
	return users
}

func handleLocationUpdate(e *Event, message []byte, useState bool) {
	update := &models.LocationUpdatePayload{}
	err := json.Unmarshal([]byte(e.Payload), update)
	if err != nil {
		log.Fatalf("Unable to unmarshal location update payload %e\n", err)
	}
	if update.Uid == "" {
		log.Fatalln("unable to find uid", update.Callsign)
	}
	if useState {
		updateLocationUsingState(update.Location, update.Uid)
	} else {
		collections.UsersCollectionManager.UpdateLocation(context.Background(), update.Location, update.Uid)
	}
	connections := []*gws.Conn{}
	connMap.Range(func(key *gws.Conn, value models.HuntUser) bool {
		connections = append(connections, key)
		return true
	})
	Broadcast(connections, gws.OpcodeText, message)
}

func updateLocationUsingState(update models.Location, uid string) {
	err := locations.UpdateLocation(update.Lat, update.Lon, update.Alt, uid)
	if err != nil {
		log.Fatalln("Unable to update client's location", err)
	}
}
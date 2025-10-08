package socket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"hunt/collections"
	"hunt/constants"
	"hunt/structs"
	"hunt/utils"
	"net"
	"net/http"
	"os"
)

func handleUpdateCallsign(event Event, c *Client) error {
	requestBody := structs.UpdateCallsignBody{}
	err := json.Unmarshal([]byte(event.Payload), &requestBody)
	if err != nil {
		return err
	}
	if requestBody.Uid == "" && requestBody.NewCallsign != "" {
		newUid := uuid.NewString()
		newUser := structs.HuntUser{
			Uid:      newUid,
			Callsign: requestBody.NewCallsign,
			Active:   true,
		}
		err = collections.LocationCollection.AddUser(c.ctx, newUser)
		if err != nil {
			return err
		}
		c.writeChan <- Event{
			Type:    NewUid,
			Payload: newUid,
		}
		return nil
	}
	if requestBody.NewCallsign == "" {
		c.writeChan <- Event{
			Type:    Error,
			Payload: "Invalid callsign",
		}
		return constants.ErrInvalidCallsign
	}
	err = collections.LocationCollection.SetCallsign(c.ctx, requestBody.Uid, requestBody.NewCallsign)
	if err != nil {
		return err
	}
	return nil
}

func handleNewGroup(event Event, c *Client) error {
	unit := structs.Unit{
		Name: event.Payload,
	}
	err := collections.GroupCollection.InsertUnit(c.ctx, unit)
	if err != nil {
		return err
	}
	updateEvent := Event{
		Type:    UnitUpdate,
		Payload: event.Payload,
	}
	err = c.manager.BroadcastMessage(updateEvent)
	if err != nil {
		return err
	}
	return nil
}

func handleTargetUpdate(event Event, c *Client) error {
	update := &structs.TargetUpdatePayload{}
	err := json.Unmarshal([]byte(event.Payload), update)
	if err != nil {
		return err
	}
	err = collections.LocationCollection.UpdateTarget(c.ctx, update.Uid, update.TargetUid, update.TargetName)
	if err != nil {
		return err
	}
	return nil
}

func handleCommandNodeUpdate(event Event, c *Client) error {
	err := c.manager.BroadcastMessage(event)
	if err != nil {
		return err
	}
	return nil
}

func handleLocationUpdate(event Event, c *Client) error {
	update := structs.HuntUser{}
	err := json.Unmarshal([]byte(event.Payload), &update)
	update.Type = "friendly"
	if err != nil {
		return err
	}
	if update.Uid == "" {
		newUid := uuid.NewString()
		connMap[c.conn] = newUid
		update.Uid = newUid
		m := Event{
			Type:    NewUid,
			Payload: newUid,
		}
		c.writeChan <- m
		return nil
	}
	update.Active = true
	err = collections.LocationCollection.UpdateUser(c.ctx, update)
	if err != nil {
		return err
	}
	err = c.manager.BroadcastMessage(event)
	if err != nil {
		return err
	}
	return nil
}

func handleInit(event Event, c *Client) error {
	incoming := structs.InitMsg{}
	err := json.Unmarshal([]byte(event.Payload), &incoming)
	if err != nil {
		return err
	}
	if incoming.Callsign == "" {
		c.writeChan <- Event{
			Type:    Error,
			Payload: "Invalid callsign",
		}
	}
	if incoming.Uid != "" {
		connMap[c.conn] = incoming.Uid
	} else {
		newUid := uuid.NewString()
		connMap[c.conn] = newUid
		c.writeChan <- Event{
			Type:    NewUid,
			Payload: newUid,
		}
	}
	result, err := collections.CommandNodeCollection.FindAll(c.ctx)
	if err != nil {
		return err
	}
	nodesToString, err := json.Marshal(result)
	if err != nil {
		return err
	}
	commandNodes := Event{
		Type:    CommandNodes,
		Payload: string(nodesToString),
	}
	c.writeChan <- commandNodes
	huntUsers, err := collections.LocationCollection.FindAllActive(c.ctx)
	if err != nil {
		return err
	}
	locationsToString, err := json.Marshal(huntUsers)
	if err != nil {
		return err
	}
	locationMsg := Event{
		Type:    Locations,
		Payload: string(locationsToString),
	}
	c.writeChan <- locationMsg
	groups, err := collections.GroupCollection.FindAll(c.ctx)
	if err != nil {
		return err
	}
	unitsToString, err := json.Marshal(groups)
	if err != nil {
		return err
	}
	unitMsg := Event{
		Type:    Units,
		Payload: string(unitsToString),
	}
	c.writeChan <- unitMsg
	return nil
}

func handleNodeStatusUpdate(event Event, c *Client) error {
	statusUpdate := structs.NodeStatusUpdate{}
	err := json.Unmarshal([]byte(event.Payload), &statusUpdate)
	if err != nil {
		return err
	}
	err = collections.CommandNodeCollection.SetStatus(c.ctx, statusUpdate.Uid, statusUpdate.Status)
	if err != nil {
		return err
	}
	return nil
}

func handleNewCommandNode(event Event, c *Client) error {
	ditHost := os.Getenv("DIT_HOST")
	ditPort := os.Getenv("DIT_CLUSTER_PORT")
	hostPort := net.JoinHostPort(ditHost, ditPort)
	reqUrl := fmt.Sprintf("http://%s/v1/createObject", hostPort)
	liveNode := structs.LiveNode{}
	err := json.Unmarshal([]byte(event.Payload), &liveNode)
	if err != nil {
		return err
	}
	b, err := json.Marshal(liveNode)
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	resp, err := utils.SendApiRequest(c.ctx, http.MethodPut, reqUrl, bytes.NewBuffer(b), headers)
	if err != nil {
		return err
	}
	resBody := utils.Response{}
	err = json.Unmarshal(resp, &resBody)
	if err != nil {
		return err
	}
	if resBody.Status == "Success" {
		msg := Event{
			Type:    NewCommandNode,
			Payload: "Success",
		}
		c.writeChan <- msg
	} else {
		msg := Event{
			Type:    Error,
			Payload: "Invalid command",
		}
		c.writeChan <- msg
	}
	return nil
}

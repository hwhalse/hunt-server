package socket

import "time"

type Event struct {
	Type    int    `json:"type"`
	Payload string `json:"payload"`
	Time time.Time `json:"time"`
}

const (
	Init                    = 1
	InitResponse			= 2
	CommandNodes            = 3
	Locations               = 4
	Units                   = 5
	LocationUpdate          = 6
	CommandNodeUpdate       = 7
	TargetUpdate            = 8
	NewCommandNode          = 9
	NewGroup                = 10
	UnitUpdate              = 11
	Disconnect              = 12
	UpdateCallsign          = 13
	NewUid                  = 14
	CommandNodeDelete       = 15
	CommandNodeStatusUpdate = 16
	Error                   = 99
)

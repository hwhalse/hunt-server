package socket

type Event struct {
	Type    int    `json:"type"`
	Payload string `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	Init                    = 1
	CommandNodes            = 2
	Locations               = 3
	Units                   = 4
	LocationUpdate          = 5
	CommandNodeUpdate       = 6
	TargetUpdate            = 7
	NewCommandNode          = 8
	NewGroup                = 9
	UnitUpdate              = 10
	Disconnect              = 11
	UpdateCallsign          = 12
	NewUid                  = 13
	CommandNodeDelete       = 14
	CommandNodeStatusUpdate = 15
	Error                   = 99
)

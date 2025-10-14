package models

type CommandNode struct {
	Uid       string   `json:"uid"`
	Callsign  string   `json:"callsign"`
	Type      string   `json:"type"`
	Location  Location `json:"location"`
	Model     string   `json:"model"`
	ModelType string   `json:"model_type"`
	Status    int      `json:"status"`
	CotType   string   `json:"cot_type"`
}

type LiveNode struct {
	Source   string   `json:"source"`
	Type     string   `json:"type"`
	CotType  string   `json:"cot_type"`
	Callsign string   `json:"callsign"`
	Location Location `json:"location"`
}

type NodeStatusUpdate struct {
	Uid    string `json:"uid"`
	Status int    `json:"status"`
}

type UpdateCallsignBody struct {
	Uid         string `json:"uid"`
	NewCallsign string `json:"newCallsign"`
}

type HuntUser struct {
	Callsign   string   `json:"callsign"`
	Type       string   `json:"type"`
	Location   Location `bson:"location" json:"location"`
	Active     bool     `json:"active"`
	TargetName string   `json:"targetName"`
	TargetUid  string   `json:"targetUid"`
	Unit       string   `json:"unit"`
	Uid        string   `json:"uid"`
}

type TargetUpdatePayload struct {
	Uid        string `json:"userUid"`
	TargetName string `json:"targetName"`
	TargetUid  string `json:"targetUid"`
}

type LocationUpdatePayload struct {
	Callsign string      `json:"callsign,omitempty"`
	Type     string      `json:"type,omitempty"`
	Location Location `json:"location,omitempty"`
	Unit     string      `json:"unit,omitempty"`
	Active   bool        `json:"active,omitempty"`
	Uid      string      `json:"uid"`
}

type InitMsg struct {
	Uid      string `json:"uid"`
	Callsign string `json:"callsign"`
}

type InitResponse struct {
	Groups      []Unit `json:"groups"`
	Users []HuntUser `json:"users"`
	CommandNodes []CommandNode `json:"commandNode"`
}

type Location struct {
	Lat float64 `bson:"lat" json:"lat"`
	Lon float64 `bson:"lon" json:"lon"`
	Alt float64 `bson:"alt" json:"alt"`
}

type Unit struct {
	Name string `json:"name"`
}

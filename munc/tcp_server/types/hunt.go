package types

type HuntUser struct {
	Callsign   string   `json:"callsign"`
	Type       string   `json:"type"`
	Location   Location `json:"location"`
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

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Alt float64 `json:"alt"`
}

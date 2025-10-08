package types

type Target struct {
	Uid       string   `json:"uid"`
	Callsign  string   `json:"callsign"`
	Type      string   `json:"type"`
	Location  Location `json:"location"`
	Model     string   `json:"model"`
	ModelType string   `json:"model_type"`
	Status    int      `json:"status"`
	CotType   string   `json:"cot_type"`
}

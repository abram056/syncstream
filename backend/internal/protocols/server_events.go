package protocols

type RoomStateEvent struct {
	Type              string  `json:"type"`
	RoomID            string  `json:"roomId"`
	DisplayName       string  `json:"displayName"`
	MediaURL          string  `json:"mediaUrl"`
	IsPlaying         bool    `json:"isPlaying"`
	Position          float64 `json:"position"`
	NumOfParticipants int     `json:"numOfParticipants"`
}

type SyncStateEvent struct {
	Type        string  `json:"type"`
	RoomID      string  `json:"roomId"`
	IsPlaying   bool    `json:"isPlaying"`
	Position    float64 `json:"position"`
	InitiatedBy string  `json:"initiatedBy"`
}

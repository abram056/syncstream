package protocols

type JoinRoomEvent struct {
	Type        string `json:"type"`
	RoomID      string `json:"roomId"`
	DisplayName string `json:"displayName"`
}

type LeaveRoomEvent struct {
	Type   string `json:"type"`
	RoomID string `json:"roomId"`
}

type SeekEvent struct {
	Type     string  `json:"type"`
	Position float64 `json:"position"`
}

type PlayEvent struct {
	Type     string  `json:"type"`
	Position float64 `json:"position"`
}

type PauseEvent struct {
	Type     string  `json:"type"`
	Position float64 `json:"position"`
}

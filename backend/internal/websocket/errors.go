package websocket

import "errors"

var (
	ErrInvalidEvent = errors.New("invalid event")
	ErrUnknownEvent = errors.New("unknown event type")
)

package websocket

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket connection to a room.
type Client struct {
	Hub           *Hub
	Conn          *websocket.Conn
	Send          chan []byte
	ParticipantID string
	DisplayName   string
	Joined        bool
}

// NewClient creates a new client connection.
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		Hub:           hub,
		Conn:          conn,
		Send:          make(chan []byte, 256),
		ParticipantID: generateParticipantID(),
		Joined:        false,
	}
}

// ReadPump reads messages from the WebSocket connection.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}

		var evt map[string]interface{}
		if err := json.Unmarshal(message, &evt); err != nil {
			log.Printf("failed to decode event: %v", err)
			continue
		}

		if err := c.handleEvent(evt); err != nil {
			log.Printf("event handler error: %v", err)
		}
	}
}

// WritePump writes messages to the WebSocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// hub closed the send channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleEvent processes an incoming event from the client.
func (c *Client) handleEvent(evt map[string]interface{}) error {
	eventType, ok := evt["type"].(string)
	if !ok {
		return ErrInvalidEvent
	}

	switch eventType {
	case "join_room":
		if !c.Joined {
			if err := c.Hub.HandleJoinRoom(c, evt); err != nil {
				return err
			}
			c.Joined = true
		}

	case "ping":
		pong := map[string]interface{}{"type": "pong"}
		if msg, err := json.Marshal(pong); err == nil {
			c.Send <- msg
		}

	default:
		return ErrUnknownEvent
	}

	return nil
}

// generateParticipantID creates a unique participant ID.
func generateParticipantID() string {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return time.Now().Format("usr20060102150405")
	}
	return "usr-" + hex.EncodeToString(buf)
}

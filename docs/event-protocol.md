## Client → Server Events

### Join Room

Sent when a user joins an existing room.

```
{  "type": "join_room",  "room_id": "rm123",  "display_name": "Alice"}
```

---

### Play

Sent when a user starts playback.

```
{  "type": "play",  "position": 142.0}
```

---

### Pause

Sent when a user pauses playback.

```
{  "type": "pause",  "position": 142.0}
```

---

### Seek

Sent when a user jumps to another timestamp.

```
{  "type": "seek",  "position": 310.5}
```

---

### Ping

Heartbeat message used to keep the connection alive and measure latency.

```
{  "type": "ping"}
```

---

# Server → Client Events

### Room Joined

Confirms successful room join.

```
{  "type": "room_joined",  "room_id": "rm123"}
```

---

### Room State

Sent:

- after joining a room
- after reconnecting
- when a client requests a full state refresh

Provides the authoritative room state.

```
{  "type": "room_state",  "room_id": "rm123", media_url: "https://youtube.com/watch?v=abc123", "is_playing": true,  "position": 142.0,  "participants": 3}
```

---

### Sync State

Broadcast whenever playback changes.

```
{  "type": "sync_state",  "room_id": "rm123",  "is_playing": false,  "position": 142.0,  "initiated_by": "usr123"}
```

Examples:

- play
- pause
- seek

This is the primary synchronization event.

---

### User Joined

Broadcast when a new participant enters the room.

```
{  "type": "user_joined",  "user_id": "usr123",  "display_name": "Alice"}
```

---

### User Left

Broadcast when a participant disconnects or leaves.

```
{  "type": "user_left",  "user_id": "usr123"}
```

---

### Room Status Transitions

The server maintains an authoritative `status` for each room and emits `room_state` updates when the status or membership changes.

- `waiting` → `active`: emitted when the first participant joins the room. Clients should expect to receive a `room_state` or `user_joined` event after a successful join.
- `active` → `idle`: emitted when the last participant leaves. The server will broadcast `user_left` events as participants disconnect; when participant count reaches zero the room `status` will be `idle`.

Cleanup behavior:

- The backend runs a periodic cleanup job that deletes rooms in the `idle` state which have been inactive past a retention threshold (MVP default: 1 hour).
- If a room is deleted while a client tries to connect, the server will respond with an error (e.g. `room not found`) or close the websocket connection. Clients should handle `error` events and `404` responses gracefully and re-create a room if appropriate.


### Pong

Response to a ping.

```
{  "type": "pong"}
```

---

### Error

Sent when a request cannot be fulfilled.

```
{  "type": "error",  "message": "Room not found"}
```

Examples:

```
{  "type": "error",  "message": "Room is full"}
```

```
{  "type": "error",  "message": "Invalid playback position"}
```

---

# Protocol Summary

### Client → Server

```
join_room, play, pause, seek, ping
```

### Server → Client

```
room_joined, room_state, sync_state, user_joined, user_left, pong, error
```

---

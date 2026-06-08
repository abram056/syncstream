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

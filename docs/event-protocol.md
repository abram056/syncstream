## Client → Server Events

### Join Room

Sent when a user joins an existing room. May include a `reconnect_token` to restore a previous session.

```
{  "type": "join_room",  "room_id": "rm123",  "display_name": "Alice"}
```

Reconnect variant:

```
{  "type": "join_room",  "room_id": "rm123",  "display_name": "Alice",  "reconnect_token": "abc123def456"}
```

---

### Leave Room

Sent when a user intentionally leaves the room. This permanently removes the participant and broadcasts `user_left`, unlike a WebSocket disconnect which is treated as temporary.

```
{  "type": "leave_room",  "room_id": "rm123"}
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

# Server → Client Events

### Room Joined

Confirms successful room join. Includes a `reconnect_token` for future reconnection.

```
{  "type": "room_joined",  "room_id": "rm123",  "reconnect_token": "abc123def456"}
```

---

### Room State

Sent:

- after joining a room
- after reconnecting
- when a client requests a full state refresh

Provides the authoritative room state. Position is calculated using `EffectivePosition()`
which accounts for elapsed time when playback is active.

```
{  "type": "room_state",  "room_id": "rm123",  "media_url": "https://youtube.com/watch?v=abc123",  "is_playing": true,  "position": 142.0,  "participants": 3}
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

### User Reconnected

Broadcast when a previously disconnected participant reconnects (instead of `user_joined`).

```
{  "type": "user_reconnected",  "user_id": "usr123",  "display_name": "Alice"}
```

---

### User Disconnected

Broadcast when a WebSocket connection drops. The participant remains in the room
and can reconnect within the grace period (default: 5 minutes).

```
{  "type": "user_disconnected",  "user_id": "usr123",  "display_name": "Alice"}
```

---

### User Left

Broadcast when a participant permanently leaves (explicit `leave_room` or cleanup expiration).

```
{  "type": "user_left",  "user_id": "usr123"}
```

---

### Room Status Transitions

The server maintains an authoritative `status` for each room and emits `room_state` updates when the status or membership changes.

- `waiting` → `active`: emitted when the first participant joins the room.
- `active` → `idle`: emitted when the last participant is removed (all participants expired or left).
  The room is retained temporarily and cleaned up if it stays idle.

Cleanup behavior:

- Disconnected participants are retained for a grace period (default: 5 minutes).
- After the grace period expires, the participant is removed and `user_left` is broadcast.
- Rooms idle longer than the retention threshold (default: 30 minutes) are deleted.
- The server runs periodic cleanup at a configurable interval (default: 10 minutes).

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
join_room, leave_room, play, pause, seek
```

### Server → Client

```
room_joined, room_state, sync_state, user_joined, user_reconnected, user_disconnected, user_left, error
```

---

# Transport-Level Heartbeat

Keepalive is handled by WebSocket Ping/Pong frames (not application-level messages):

- Server sends Ping frames every 54 seconds.
- Client responds with Pong frames automatically.
- Read deadline refreshes on each Pong (60 second timeout).
- Dead connections are detected and cleaned up automatically.

This replaces the previous application-level `ping`/`pong` events.

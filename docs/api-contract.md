
# Virtual TV Room API Contract v1

## HTTP Endpoints

### Create Room

```
POST /api/v1/rooms
```

Request:

```
{   media_url: "https://example.com/video.mp4" }
```

Response:

```
{  "room_id": "rm123"}
```

Status:

```
201 Created
```

---

### Get Room Information

```
GET /api/v1/rooms/{room_id}
```

Response:

```
{  "room_id": "rm123",  "status": "active",  "participants": 3}
```

Status:

```
200 OK
404 Not Found
```

---

### Room Status & Cleanup

Rooms expose a `status` field with these possible values:

- `waiting` — room created but no active participants yet.
- `active` — one or more participants are connected and the room is actively used.
- `idle` — no participants currently connected; the room is kept temporarily.

Behavior:

- When the first participant successfully joins a room it transitions `waiting` → `active`.
- When the last participant leaves it transitions `active` → `idle`.
- Idle rooms are retained temporarily and periodically cleaned up by the backend (see cleanup policy below).

Cleanup policy:

- Disconnected participants are retained for a grace period (default: 5 minutes).
- After the grace period ends, expired participants are removed and `user_left` is broadcast.
- The server removes idle rooms that have been inactive longer than the configured retention threshold (default: 30 minutes).
- Cleanup runs in a background task at a configurable interval (default: 10 minutes).
- Deleted rooms will return `404 Not Found` for subsequent `GET /api/v1/rooms/{room_id}` requests.


### Health Check

Useful later for deployment.

```
GET /health
```

Response:

```
{  "status": "ok"}
```

---

# WebSocket Contract

This is where the real action happens.

---

## Connect To Room

```
ws://server/api/v1/rooms/{room_id}/ws
```

Example:

```
ws://localhost:8080/api/v1/rooms/rm123/ws
```

After connecting:

Client sends:

```
{  "type": "join_room",  "room_id": "rm123",  "display_name": "Alice"}
```

---

# Event Flow

Once connected, all communication follows the event protocol you already defined.

### Client → Server

```
join_room (with optional reconnect_token)
leave_room
play
pause
seek
```

---

### Server → Client

```
room_joined (includes reconnect_token)
room_state (authoritative position)
sync_state
user_joined
user_reconnected
user_disconnected
user_left
error
```

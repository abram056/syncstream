
# Virtual TV Room API Contract v1

## HTTP Endpoints

### Create Room

```
POST /api/v1/rooms
```

Request:

```
{}
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
join_room
play
pause
seek
ping
```

---

### Server → Client

```
room_joined
room_state
sync_state
user_joined
user_left
pong
error
```

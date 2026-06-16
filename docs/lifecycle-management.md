# Lifecycle Management

## Participant Lifecycle

### States
- **Connected**: Participant has an active WebSocket connection.
- **Disconnected**: Participant's WebSocket closed, but the participant record remains
  in the room with `Connected=false` and `LastSeen` updated.
- **Expired**: Disconnected for longer than the grace period (default: 5 minutes),
  the participant is permanently removed.

### Transitions
```
Connected ──WebSocket close──► Disconnected ──grace period──► Expired (removed)
    ▲                              │
    └─────Reconnect (within ───────┘
           grace period)
```

### Reconnect Flow
1. Client connects WebSocket, sends `join_room` with `reconnect_token`
2. Server validates token via `Manager.ReconnectParticipant`
3. On success: participant marked `Connected=true`, `user_reconnected` broadcast
4. On failure: server falls back to creating a new participant

### Disconnect vs Leave
| Event | Trigger | Effect |
|-------|---------|--------|
| `user_disconnected` | WebSocket close or network loss | Participant remains in room |
| `user_left` | Cleanup expiration or explicit `leave_room` | Participant removed permanently |

## Playback State

### Authoritative Position
When `IsPlaying=true`, the effective position is calculated as:
```
effective_position = stored_position + time_since_last_update
```

This ensures:
- Reconnecting users receive the current playback position
- Late joiners see accurate timing
- Paused playback returns the exact stored position

## Room Lifecycle

### States
- **Waiting**: Created, no participants
- **Active**: At least one participant
- **Idle**: All participants removed (room retained temporarily)
- **Deleted**: Idle longer than timeout (default: 30 minutes)

### Cleanup Services
Two background goroutines run at configurable intervals (default: 10 min):

1. **Participant Expiration**: Removes disconnected participants past grace period.
   Broadcasts `user_left` for each expired participant via the room's hub.

2. **Idle Room Cleanup**: Removes rooms idle longer than timeout.
   Also removes the corresponding hub via `HubRegistry.Remove()`.

## Hub Lifecycle

### HubRegistry
Centralizes WebSocket hub management:
- `GetOrCreate(roomID, manager)`: Returns existing hub or creates a new one
- `Get(roomID)`: Returns hub if it exists
- `Remove(roomID)`: Stops and removes a hub (called when room is deleted)
- `List()`: Returns all active hubs

### Hub Shutdown
`Hub.Stop()` signals the `Run()` goroutine to exit. On shutdown:
1. All client Send channels are closed
2. All WebSocket connections are closed
3. The goroutine returns, releasing its resources

## Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `ParticipantGracePeriod` | 5 min | How long a disconnected participant is retained |
| `RoomIdleTimeout` | 30 min | How long an idle room is retained before cleanup |
| `CleanupInterval` | 10 min | How often the cleanup loop runs |

## Concurrency

Room state is protected by an internal `sync.RWMutex`. All mutations via
`Manager` methods acquire a write lock; read-only operations use a read lock.
The `RoomStore` map is separately protected by its own `sync.RWMutex`.

## Event Protocol Additions

### Client → Server
- `leave_room`: Explicit room departure (triggers `user_left`)
- `join_room` now accepts optional `reconnect_token` field

### Server → Client
- `user_disconnected`: Temporary WebSocket loss
- `user_reconnected`: Participant restored their session
- `join_room` response now includes `reconnect_token` for future reconnection

### Removed
- Application-level `ping`/`pong` events (WebSocket Ping/Pong handles keepalive)

## Known Limitations

1. **Race on room deletion**: If a room transitions from idle to active between the
   cleanup check and deletion, the room could be deleted while in use. Mitigation:
   verify room presence before WebSocket upgrade in `ServeWS`.

2. **Token security**: Reconnect tokens are random hex strings without expiration.
   In production, tokens should have an expiry and be tied to a specific client IP
   or session.

3. **Single-node in-memory**: All state is in-memory. Server restart loses all rooms
   and participants. For production, add persistent storage.

4. **Broadcast after hub stop**: If a participant expires and its hub is stopped
   concurrently, the `user_left` broadcast may be lost. Currently handled by checking
   hub existence before broadcasting.

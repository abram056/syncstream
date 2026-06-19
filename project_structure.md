````
# Project Structure

```
watchroom/
│
├── backend/
│   │
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   │
│   ├── internal/
│   │   │
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   │   ├── create_room.go
│   │   │   │   ├── get_room.go
│   │   │   │   └── health.go
│   │   │   │
│   │   │   └── router.go
│   │   │
│   │   ├── websocket/
│   │   │   ├── handler.go
│   │   │   ├── client.go
│   │   │   └── hub.go
│   │   │
│   │   ├── room/
│   │   │   ├── room.go
│   │   │   ├── manager.go
│   │   │   └── repository.go
│   │   │
│   │   ├── playback/
│   │   │   ├── state.go
│   │   │   └── service.go
│   │   │
│   │   ├── protocol/
│   │   │   ├── client_events.go
│   │   │   ├── server_events.go
│   │   │   └── event_types.go
│   │   │
│   │   ├── models/
│   │   │   ├── room.go
│   │   │   ├── participant.go
│   │   │   ├── media.go
│   │   │   ├── playback_state.go
│   │   │   └── room_status.go
│   │   │
│   │   └── storage/
│   │       └── memory/
│   │           └── room_store.go
│   │
│   ├── configs/
│   │   └── config.go
│   │
│   ├── tests/
│   │   ├── integration/
│   │   └── unit/
│   │
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   │
│   ├── src/
│   │   │
│   │   ├── components/
│   │   │   ├── VideoPlayer.tsx
│   │   │   ├── ParticipantList.tsx
│   │   │   └── PlaybackControls.tsx
│   │   │
│   │   ├── pages/
│   │   │   ├── Home.tsx
│   │   │   └── Room.tsx
│   │   │
│   │   ├── services/
│   │   │   ├── api.ts
│   │   │   └── websocket.ts
│   │   │
│   │   ├── hooks/
│   │   │   └── useRoom.ts
│   │   │
│   │   ├── types/
│   │   │   └── protocol.ts
│   │   │
│   │   ├── App.tsx
│   │   └── main.tsx
│   │
│   ├── package.json
│   └── vite.config.ts
│
├── docs/
│   │
│   ├── architecture/
│   │   ├── architecture.md
│   │   ├── component-diagram.png
│   │   └── deployment.md
│   │
│   ├── protocols/
│   │   ├── api-contract.md
│   │   └── websocket-events.md
│   │
│   ├── domain/
│   │   ├── domain-model.md
│   │   └── class-diagram.png
│   │
│   ├── diagrams/
│   │   ├── join-room-sequence.png
│   │   ├── playback-sequence.png
│   │   └── reconnect-sequence.png
│   │
│   └── decisions/
│       ├── 001-server-authoritative-sync.md
│       ├── 002-no-media-hosting.md
│       ├── 003-websocket-communication.md
│       └── 004-remote-media-url.md
│
├── .gitignore
├── README.md
└── LICENSE
```

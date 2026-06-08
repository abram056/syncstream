## 1. Project Overview

Virtual TV Room is a realtime synchronized watch-party platform that allows multiple users to remotely watch video content together inside shared rooms.

The platform is primarily responsible for:

* coordinating playback state between participants
* managing shared rooms
* synchronizing playback events in realtime
* maintaining participant presence and room state

The platform is *not* intended to function as a full media hosting or streaming provider during the MVP phase.

---

# 2. Architectural Goals

The system architecture prioritizes the following:

* modularity
* scalability
* maintainability
* low coupling between components
* realtime responsiveness
* clean separation of responsibilities

The MVP will initially deploy as a modular monolithic backend while maintaining internal boundaries that can later evolve into distributed services if required.

---

# 3. MVP Scope

## Included Features

* room creation
* room joining via room code/link
* synchronized playback
* play/pause/seek synchronization
* realtime participant presence
* room state management
* temporary session handling
* WebSocket-based realtime communication

## Excluded Features

The following features are intentionally excluded from the MVP:

* permanent video hosting
* DRM-protected media redistribution
* automatic video downloading
* voice chat
* video chat
* recommendation systems
* persistent social features
* distributed multi-region deployment
* advanced analytics
* mobile applications
* media transcoding pipelines

---

# 4. System Responsibility Boundaries

## Responsibilities of the Platform

The backend system is responsible for:

* maintaining authoritative room state
* coordinating playback synchronization
* broadcasting realtime events
* managing room lifecycle
* tracking connected participants
* handling reconnect/disconnect logic

## Non-Responsibilities

The backend system is *not* responsible for:

* long-term copyrighted media storage
* unauthorized media redistribution
* content acquisition from external platforms
* media transcoding
* CDN/media delivery infrastructure

---

# 5. Media Handling Model

The system follows a synchronization-first architecture rather than a media-hosting architecture.

## Media Sources

Video content may originate from:

* user-provided URLs
* legally embeddable/public video sources
* temporary uploads (future possibility)
* local client-side media sources

## Backend Media Policy

The backend:

* does not permanently store copyrighted content
* does not download media from external providers
* does not redistribute DRM-protected content
* primarily stores metadata and playback state

The server coordinates playback state only.

---

# 6. Scale Assumptions (MVP)

Initial target scale for MVP:

* 50 concurrent users
* 10 active rooms
* 2–8 users per room
* single backend instance
* single deployment region
* in-memory synchronization state

These assumptions allow the architecture to remain relatively simple while still exercising realtime synchronization and state-management concerns.

---

# 7. High-Level System Architecture

## Primary Architectural Style

The MVP uses a modular monolithic backend architecture.

Internally, the backend is separated into logical components/modules.

Example logical structure:

* API Layer
* WebSocket Gateway
* Room Manager
* Playback Manager
* Synchronization Manager
* Presence Manager
* Session Store

These boundaries are intended to reduce coupling and allow future extraction into services if scaling demands increase.

---

# 8. Communication Model

## HTTP Responsibilities

HTTP endpoints are used for:

* room creation
* room lookup
* authentication (future)
* metadata retrieval
* initial session negotiation

## WebSocket Responsibilities

WebSockets are used for:

* realtime playback synchronization
* room events
* participant presence
* event broadcasting
* synchronization heartbeats

Realtime communication is expected to remain persistent for the duration of a room session.

---

# 9. Source of Truth Model

The backend server acts as the authoritative source of truth for room playback state.

Clients:

* send playback intents/events
* receive authoritative synchronization updates from the server

The server:

* validates playback actions
* updates official room state
* broadcasts synchronized state to participants

This prevents state divergence between clients.

---

# 10. Playback Synchronization Principles

## Synchronization Goals

The synchronization system aims to:

* minimize playback drift
* maintain shared viewing consistency
* recover gracefully from reconnects
* support late joiners

## Initial Synchronization Assumptions

* server-authoritative playback timing
* periodic synchronization heartbeats
* acceptable playback drift tolerance: ±500ms
* reconnecting users receive latest room state
* joining users synchronize to current playback timestamp

Future improvements may include:

* adaptive drift correction
* latency compensation
* clock offset estimation
* predictive synchronization

---

# 11. Room Lifecycle

A room may transition through several states:

1. Created
2. WaitingForParticipants
3. Active
4. Idle
5. Expired
6. Destroyed

## Lifecycle Notes

* rooms are created on demand
* inactive rooms may expire automatically
* disconnects do not immediately destroy rooms
* room state persists temporarily during short disconnections

---

# 12. Presence Management

The system maintains participant presence within rooms.

Tracked information may include:

* connection status
* last activity timestamp
* current playback synchronization status
* participant role (future)

Presence updates are distributed through realtime events.

---

# 13. Failure Assumptions

The architecture assumes the possibility of:

* temporary network instability
* delayed packets
* client disconnects
* stale synchronization state
* simultaneous conflicting playback actions

The backend should attempt graceful recovery whenever possible.

Examples:

* reconnect users to latest room state
* resolve playback conflicts via server authority
* remove stale connections after timeout

---

# 14. Data Persistence Strategy

## MVP Persistence Model

Initial MVP persistence may remain minimal.

Potential persisted data:

* room metadata
* temporary session state
* audit/logging information

Realtime synchronization state may remain in memory during MVP development.

Future scaling phases may introduce:

* Redis
* persistent databases
* distributed session stores

---

# 15. Security Considerations

Initial security concerns include:

* room access validation
* event validation
* connection authentication (future)
* rate limiting
* invalid event rejection

Advanced security concerns are deferred beyond MVP scope.

---

# 16. Future Scalability Considerations

The architecture should remain extensible toward:

* horizontal scaling
* distributed synchronization
* external state stores
* pub/sub systems
* load balancing
* service decomposition
* regional deployments

These are considered future evolutions rather than MVP requirements.

---

# 17. Initial Technical Direction (Tentative)

Potential technology direction:

Backend:

* Go
* WebSockets
* modular monolithic architecture

Frontend:

* web-based client application

Infrastructure:

* single server deployment
* in-memory synchronization state

Technology choices remain flexible during early architecture exploration.

---

# 18. Phase 0 Deliverables

The following deliverables are expected from Phase 0:

* high-level architecture diagram
* component definitions
* playback synchronization model
* room lifecycle definitions
* event protocol definitions
* MVP scope definition
* UML diagrams

  * component diagrams
  * sequence diagrams

These deliverables establish the foundational system design before implementation begins.

// --- REST API types ---

export interface CreateRoomRequest {
  media_url: string
  title?: string
}

export interface CreateRoomResponse {
  room_id: string
}

export interface GetRoomResponse {
  room_id: string
  status: RoomStatus
  media_url: string
  title?: string
  is_playing: boolean
  position: number
  participants: number
}

export type RoomStatus = 'waiting' | 'active' | 'idle'

// --- WebSocket client → server event types ---

export type ClientEvent =
  | JoinRoomEvent
  | LeaveRoomEvent
  | PlayEvent
  | PauseEvent
  | SeekEvent

export interface JoinRoomEvent {
  type: 'join_room'
  room_id: string
  display_name: string
  reconnect_token?: string
  participant_id?: string
}

export interface LeaveRoomEvent {
  type: 'leave_room'
  room_id: string
}

export interface PlayEvent {
  type: 'play'
  position?: number
}

export interface PauseEvent {
  type: 'pause'
  position?: number
}

export interface SeekEvent {
  type: 'seek'
  position: number
}

// --- WebSocket server → client event types ---

export type ServerEvent =
  | RoomJoinedEvent
  | RoomStateEvent
  | SyncStateEvent
  | UserJoinedEvent
  | UserReconnectedEvent
  | UserDisconnectedEvent
  | UserLeftEvent
  | ErrorEvent

export interface RoomJoinedEvent {
  type: 'room_joined'
  roomId: string
  participantId: string
  reconnect_token: string
}

export interface RoomStateEvent {
  type: 'room_state'
  roomId: string
  status: string
  mediaUrl: string
  isPlaying: boolean
  position: number
  numOfParticipants: number
  updatedBy?: string
  updatedAt?: number
  participants?: Array<{ userId: string; displayName: string; connected: boolean }>
}

export interface SyncStateEvent {
  type: 'sync_state'
  roomId: string
  isPlaying: boolean
  position: number
  updatedBy: string
  updatedAt: number
  numOfParticipants: number
}

export interface UserJoinedEvent {
  type: 'user_joined'
  userId: string
  displayName: string
}

export interface UserReconnectedEvent {
  type: 'user_reconnected'
  userId: string
  displayName: string
}

export interface UserDisconnectedEvent {
  type: 'user_disconnected'
  userId: string
  displayName: string
}

export interface UserLeftEvent {
  type: 'user_left'
  userId: string
  displayName?: string
}

export interface ErrorEvent {
  type: 'error'
  message: string
}

// --- Application state types ---

export interface Participant {
  userId: string
  displayName: string
  connected: boolean
}

export type RoomView = 'lobby' | 'watching'

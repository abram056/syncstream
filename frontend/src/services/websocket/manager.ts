import type { ClientEvent, ServerEvent } from '../../types'

type EventHandler = (event: ServerEvent) => void
type StatusHandler = (status: ConnectionStatus) => void

export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected'

export class WebSocketManager {
  private ws: WebSocket | null = null
  private url: string = ''
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 2000
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null
  private eventHandlers: Set<EventHandler> = new Set()
  private statusHandlers: Set<StatusHandler> = new Set()
  private _status: ConnectionStatus = 'disconnected'
  private roomId: string = ''
  private shouldReconnect = false

  get status() {
    return this._status
  }

  connect(roomId: string, _displayName: string) {
    this.roomId = roomId
    this.shouldReconnect = true
    this.reconnectAttempts = 0
    this.doConnect()
  }

  private doConnect() {
    if (this.ws?.readyState === WebSocket.OPEN) return

    this.setStatus('connecting')

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    this.url = `${protocol}//${host}/api/v1/rooms/${this.roomId}/ws`

    this.ws = new WebSocket(this.url)

    this.ws.onopen = () => {
      this.reconnectAttempts = 0
      this.setStatus('connected')
    }

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as ServerEvent
        this.dispatch(data)
      } catch {
        // ignore malformed messages
      }
    }

    this.ws.onclose = () => {
      this.setStatus('disconnected')
      this.ws = null
      if (this.shouldReconnect) {
        this.scheduleReconnect()
      }
    }

    this.ws.onerror = () => {
      // onclose will fire after this
    }
  }

  private scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) return
    this.reconnectAttempts++
    const delay = this.reconnectDelay * Math.pow(1.5, this.reconnectAttempts - 1)
    this.reconnectTimer = setTimeout(() => this.doConnect(), delay)
  }

  disconnect() {
    this.shouldReconnect = false
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    this.ws?.close()
    this.ws = null
    this.setStatus('disconnected')
  }

  send(event: ClientEvent) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(event))
    }
  }

  onEvent(handler: EventHandler): () => void {
    this.eventHandlers.add(handler)
    return () => this.eventHandlers.delete(handler)
  }

  onStatusChange(handler: StatusHandler): () => void {
    this.statusHandlers.add(handler)
    return () => this.statusHandlers.delete(handler)
  }

  private dispatch(event: ServerEvent) {
    for (const handler of this.eventHandlers) {
      handler(event)
    }
  }

  private setStatus(status: ConnectionStatus) {
    this._status = status
    for (const handler of this.statusHandlers) {
      handler(status)
    }
  }
}

// Singleton instance for the application
export const wsManager = new WebSocketManager()

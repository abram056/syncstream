import { useEffect, useState, useCallback, useRef } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Button } from '../components/ui/button'
import { Lobby } from '../components/Lobby'
import { VideoPlayer } from '../components/VideoPlayer'
import { ParticipantList } from '../components/ParticipantList'
import { wsManager } from '../services/websocket/manager'
import { useRoomStore } from '../stores/roomStore'
import { usePlaybackStore } from '../stores/playbackStore'
import { useConnectionStore } from '../stores/connectionStore'
import type { ServerEvent } from '../types'
import { JoinRoomPrompt } from "../components/JoinRoomPrompt";

export function Room() {
  const { roomId } = useParams<{ roomId: string }>()
  const navigate = useNavigate()
  const [showParticipants, setShowParticipants] = useState(false)
  const [pageError, setPageError] = useState('')
  const initialSyncDoneRef = useRef(false)

  const roomIdStored = useRoomStore((s) => s.roomId)
  const mediaUrl = useRoomStore((s) => s.mediaUrl)
  const localDisplayName = useRoomStore((s) => s.localDisplayName)
  const view = useRoomStore((s) => s.view)
  const setView = useRoomStore((s) => s.setView)
  const addParticipant = useRoomStore((s) => s.addParticipant)
  const removeParticipant = useRoomStore((s) => s.removeParticipant)
  const setParticipantConnected = useRoomStore((s) => s.setParticipantConnected)
  const setRoomStatus = useRoomStore((s) => s.setRoomStatus)
  const setMediaUrl = useRoomStore((s) => s.setMediaUrl)
  const resetRuntime = useRoomStore((s) => s.resetRuntime)

  const syncState = usePlaybackStore((s) => s.syncState)
  const isPlaying = usePlaybackStore((s) => s.isPlaying)
  const resetPlayback = usePlaybackStore((s) => s.reset)

  const setConnectionStatus = useConnectionStore((s) => s.setStatus)
  const incrementReconnect = useConnectionStore((s) => s.incrementReconnect)
  const resetReconnect = useConnectionStore((s) => s.resetReconnect)

  const setLocalDisplayName = useRoomStore((s) => s.setLocalDisplayName)
  const reconnectToken = useRoomStore((s) => s.reconnectToken)
  const participantId = useRoomStore((s) => s.participantId)
  const setReconnectInfo = useRoomStore((s) => s.setReconnectInfo)

  // Redirect if we don't have room context
  useEffect(() => {
    if (!roomIdStored && !roomId) {
      navigate('/')
    }
  }, [roomIdStored, roomId, navigate])

  const effectiveRoomId = roomId || roomIdStored || ''

  // Set up WebSocket and event handlers
  useEffect(() => {
    if (!effectiveRoomId || !localDisplayName) return
    if (initialSyncDoneRef.current) return

    const handleEvent = (event: ServerEvent) => {
      switch (event.type) {
        case 'room_joined': {
          setReconnectInfo(event.participantId, event.reconnect_token)
          break
        }
        case 'room_state': {
          setRoomStatus(event.status)
          setMediaUrl(event.mediaUrl)
          syncState({
            isPlaying: event.isPlaying,
            position: event.position,
            updatedBy: event.updatedBy || '',
            updatedAt: event.updatedAt || Date.now(),
          })
          if (event.participants) {
            for (const p of event.participants) {
              addParticipant(p.userId, p.displayName)
            }
          }
          initialSyncDoneRef.current = true
          break
        }
        case 'sync_state': {
          syncState({
            isPlaying: event.isPlaying,
            position: event.position,
            updatedBy: event.updatedBy,
            updatedAt: event.updatedAt,
          })
          break
        }
        case 'user_joined': {
          addParticipant(event.userId, event.displayName)
          break
        }
        case 'user_reconnected': {
          setParticipantConnected(event.userId, true)
          break
        }
        case 'user_disconnected': {
          setParticipantConnected(event.userId, false)
          break
        }
        case 'user_left': {
          removeParticipant(event.userId)
          break
        }
        case 'error': {
          setPageError(event.message)
          break
        }
      }
    }

    const handleStatus = (status: string) => {
      setConnectionStatus(status as 'disconnected' | 'connecting' | 'connected')
      if (status === 'connected') {
        resetReconnect()
        // Send join_room
        wsManager.send({
          type: 'join_room',
          room_id: effectiveRoomId,
          display_name: localDisplayName,
          ...(reconnectToken && participantId ? {
            reconnect_token: reconnectToken,
            participant_id: participantId,
          } : {}),
        })
      } else if (status === 'disconnected') {
        incrementReconnect()
      }
    }

    const unsubEvent = wsManager.onEvent(handleEvent)
    const unsubStatus = wsManager.onStatusChange(handleStatus)

    wsManager.connect(effectiveRoomId, localDisplayName)

    return () => {
      unsubEvent()
      unsubStatus()
    }
  }, [effectiveRoomId, localDisplayName])

  const handleTogglePlay = useCallback((position?: number) => {
    if (view === 'lobby') {
      wsManager.send({
        type: 'play', position
      })
      setView('watching')
    } else {
      if (isPlaying) {
        wsManager.send({ type: 'pause', position })
      } else {
        wsManager.send({ type: 'play', position })
      }
    }
  }, [view, isPlaying, setView])

  const handleSeek = useCallback((position: number) => {
    wsManager.send({ type: 'seek', position })
  }, [])

  const handleStartWatching = useCallback(() => {
    wsManager.send({ type: 'play' })
    setView('watching')
  }, [setView])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      wsManager.disconnect()
      resetRuntime()
      resetPlayback()
      initialSyncDoneRef.current = false
    }
  }, [])

  if (pageError) {
    return (
      <div className="flex min-h-[calc(100vh-3.5rem)] flex-col items-center justify-center gap-4">
        <p className="text-red-400">{pageError}</p>
        <Button onClick={() => navigate('/')}>Go Home</Button>
      </div>
    )
  }

  if (!localDisplayName) {
    return (
      <JoinRoomPrompt
        onSubmit={(name) => {
          setLocalDisplayName(name)
        }}
      />
    )
  }

  return (
    <div className="flex min-h-[calc(100vh-3.5rem)] flex-col">
      {view === 'lobby' ? (
        <Lobby onStartWatching={handleStartWatching} roomId={effectiveRoomId} />
      ) : (
        <div className="flex flex-1 flex-col">
          {mediaUrl && (
            <VideoPlayer
              mediaUrl={mediaUrl}
              onTogglePlay={handleTogglePlay}
              onSeek={handleSeek}
              onShowParticipants={() => setShowParticipants(true)}
            />
          )}
        </div>
      )}

      <ParticipantList
        open={showParticipants}
        onClose={() => setShowParticipants(false)}
      />
    </div>
  )
}

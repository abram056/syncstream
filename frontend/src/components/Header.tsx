import { useNavigate } from 'react-router-dom'
import { Monitor, LogOut } from 'lucide-react'
import { useConnectionStore } from '../stores/connectionStore'
import { Badge } from './ui/badge'
import { Button } from './ui/button'
import { useRoomStore } from '../stores/roomStore'
import { usePlaybackStore } from '../stores/playbackStore'
import { wsManager } from '../services/websocket/manager'

export function Header() {
  const navigate = useNavigate()
  const connectionStatus = useConnectionStore((s) => s.status)
  const roomId = useRoomStore((s) => s.roomId)
  const leaveRoom = useRoomStore((s) => s.leaveRoom)
  const resetPlayback = usePlaybackStore((s) => s.reset)


  const handleLeave = () => {
    if (roomId) {
      wsManager.send({ type: 'leave_room', room_id: roomId })
    }
    wsManager.disconnect()
    leaveRoom()
    resetPlayback()
    navigate('/')
  }

  return (
    <header className="flex h-14 items-center justify-between border-b border-zinc-800 px-4">
      <button
        onClick={() => navigate('/')}
        className="flex items-center gap-2 text-lg font-semibold"
      >
        <Monitor className="h-5 w-5 text-blue-500" />
        SyncStream
      </button>

      {connectionStatus !== 'disconnected' && (
        <Badge
          variant={
            connectionStatus === 'connected' ? 'success' : 'warning'
          }
        >
          {connectionStatus === 'connected' ? 'Connected' : 'Reconnecting...'}
        </Badge>
      )}

      {roomId && (
        <Button variant="ghost" size="sm" onClick={handleLeave}>
          <LogOut className="h-4 w-4 mr-1" />
          Leave
        </Button>
      )}
    </header>
  )
}

import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { ArrowLeft, LogIn } from 'lucide-react'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Card } from '../components/ui/card'
import { getRoom } from '../services/api/client'
import { useRoomStore } from '../stores/roomStore'

export function JoinRoom() {
  const navigate = useNavigate()
  const setRoomId = useRoomStore((s) => s.setRoomId)
  const setMediaUrl = useRoomStore((s) => s.setMediaUrl)
  const setLocalDisplayName = useRoomStore((s) => s.setLocalDisplayName)
  const [roomIdInput, setRoomIdInput] = useState('')
  const [displayName, setDisplayName] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleJoin = async () => {
    if (!roomIdInput.trim()) {
      setError('Room ID is required')
      return
    }
    if (!displayName.trim()) {
      setError('Display name is required')
      return
    }

    setLoading(true)
    setError('')

    try {
      const room = await getRoom(roomIdInput.trim())
      setRoomId(room.room_id)
      setMediaUrl(room.media_url)
      setLocalDisplayName(displayName.trim())
      navigate(`/room/${room.room_id}`)
    } catch {
      setError('Room not found or server is unavailable')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-[calc(100vh-3.5rem)] flex-col items-center justify-center p-4">
      <div className="w-full max-w-md space-y-6">
        <Button variant="ghost" onClick={() => navigate('/')}>
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back
        </Button>

        <div className="space-y-2">
          <h1 className="text-2xl font-bold">Join a Room</h1>
          <p className="text-sm text-zinc-400">
            Enter the room ID and your display name
          </p>
        </div>

        <Card className="space-y-4">
          <Input
            id="room-id"
            label="Room ID"
            placeholder="e.g. rm-abc123"
            value={roomIdInput}
            onChange={(e) => setRoomIdInput(e.target.value)}
            disabled={loading}
          />
          <Input
            id="display-name"
            label="Your Name"
            placeholder="Enter your display name"
            value={displayName}
            onChange={(e) => setDisplayName(e.target.value)}
            disabled={loading}
          />

          {error && (
            <p className="text-sm text-red-400">{error}</p>
          )}

          <Button
            className="w-full"
            size="lg"
            onClick={handleJoin}
            disabled={loading}
          >
            <LogIn className="mr-2 h-5 w-5" />
            {loading ? 'Joining...' : 'Join Room'}
          </Button>
        </Card>
      </div>
    </div>
  )
}

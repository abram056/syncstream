import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { ArrowLeft, Film } from 'lucide-react'
import { Button } from '../components/ui/button'
import { Input } from '../components/ui/input'
import { Card } from '../components/ui/card'
import { createRoom } from '../services/api/client'
import { useRoomStore } from '../stores/roomStore'

export function CreateRoom() {
  const navigate = useNavigate()
  const setRoomId = useRoomStore((s) => s.setRoomId)
  const setMediaUrl = useRoomStore((s) => s.setMediaUrl)
  const setLocalDisplayName = useRoomStore((s) => s.setLocalDisplayName)
  const [mediaUrlInput, setMediaUrlInput] = useState('')
  const [displayName, setDisplayName] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleCreate = async () => {
    if (!mediaUrlInput.trim()) {
      setError('Media URL is required')
      return
    }
    if (!displayName.trim()) {
      setError('Display name is required')
      return
    }

    setLoading(true)
    setError('')

    try {
      const res = await createRoom({ media_url: mediaUrlInput.trim() })
      setRoomId(res.room_id)
      setMediaUrl(mediaUrlInput.trim())
      setLocalDisplayName(displayName.trim())
      navigate(`/room/${res.room_id}`)
    } catch {
      setError('Failed to create room. Is the server running?')
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
          <h1 className="text-2xl font-bold">Create a Room</h1>
          <p className="text-sm text-zinc-400">
            Enter a video URL and your display name to get started
          </p>
        </div>

        <Card className="space-y-4">
          <Input
            id="media-url"
            label="Video URL"
            placeholder="https://example.com/video.mp4"
            value={mediaUrlInput}
            onChange={(e) => setMediaUrlInput(e.target.value)}
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
            onClick={handleCreate}
            disabled={loading}
          >
            <Film className="mr-2 h-5 w-5" />
            {loading ? 'Creating...' : 'Create Room'}
          </Button>
        </Card>
      </div>
    </div>
  )
}

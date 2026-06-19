import { useNavigate } from 'react-router-dom'
import { Monitor, Plus, LogIn } from 'lucide-react'
import { Button } from '../components/ui/button'

export function Landing() {
  const navigate = useNavigate()

  return (
    <div className="flex min-h-[calc(100vh-3.5rem)] flex-col items-center justify-center p-4">
      <div className="text-center space-y-6 max-w-md">
        <div className="flex justify-center">
          <div className="rounded-2xl bg-blue-600/10 p-4">
            <Monitor className="h-12 w-12 text-blue-500" />
          </div>
        </div>

        <div className="space-y-2">
          <h1 className="text-3xl font-bold tracking-tight">SyncStream</h1>
          <p className="text-zinc-400">
            Watch videos together in real time.
            <br />
            Share a room, sync playback, enjoy together.
          </p>
        </div>

        <div className="flex flex-col gap-3">
          <Button size="lg" onClick={() => navigate('/create-room')}>
            <Plus className="mr-2 h-5 w-5" />
            Create Room
          </Button>
          <Button
            variant="secondary"
            size="lg"
            onClick={() => navigate('/join')}
          >
            <LogIn className="mr-2 h-5 w-5" />
            Join Room
          </Button>
        </div>
      </div>
    </div>
  )
}

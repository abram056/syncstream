import { useNavigate } from 'react-router-dom'
import { Monitor } from 'lucide-react'
import { useConnectionStore } from '../stores/connectionStore'
import { Badge } from './ui/badge'

export function Header() {
  const navigate = useNavigate()
  const connectionStatus = useConnectionStore((s) => s.status)

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
    </header>
  )
}

import { useState } from 'react'
import { Play, Users, Link as LinkIcon, Check } from 'lucide-react'
import { Button } from './ui/button'
import { Card } from './ui/card'
import { ParticipantList } from './ParticipantList'
import { useRoomStore } from '../stores/roomStore'

interface Props {
  onStartWatching: () => void
  roomId: string
}

export function Lobby({ onStartWatching, roomId }: Props) {
  const [showParticipants, setShowParticipants] = useState(false)
  const [copied, setCopied] = useState(false)
  const participants = useRoomStore((s) => s.participants)

  const copyRoomLink = () => {
    const link = `${window.location.origin}/room/${roomId}`
    navigator.clipboard.writeText(link)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="flex flex-col items-center justify-center gap-8 p-4 pt-16">
      <div className="text-center space-y-2">
        <h1 className="text-2xl font-bold">Ready to watch?</h1>
        <p className="text-zinc-400">
          Share the room link so others can join
        </p>
      </div>

      <Card className="w-full max-w-md space-y-4">
        <div className="flex items-center gap-2 rounded-lg bg-zinc-950 px-3 py-2 text-sm text-zinc-300">
          <span className="flex-1 truncate">
            {window.location.origin}/room/{roomId}
          </span>
          <Button variant="ghost" size="sm" onClick={copyRoomLink}>
            {copied ? (
              <Check className="h-4 w-4 text-emerald-500" />
            ) : (
              <LinkIcon className="h-4 w-4" />
            )}
          </Button>
        </div>

        <div className="flex items-center justify-between">
          <Button variant="ghost" onClick={() => setShowParticipants(true)}>
            <Users className="mr-2 h-4 w-4" />
            {participants.length} participant{participants.length !== 1 ? 's' : ''}
          </Button>
        </div>

        <Button className="w-full" size="lg" onClick={onStartWatching}>
          <Play className="mr-2 h-5 w-5" />
          Start Watching
        </Button>
      </Card>

      <ParticipantList
        open={showParticipants}
        onClose={() => setShowParticipants(false)}
      />
    </div>
  )
}

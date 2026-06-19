import { X, Users } from 'lucide-react'
import { useRoomStore } from '../stores/roomStore'
import { Button } from './ui/button'

interface Props {
  open: boolean
  onClose: () => void
}

export function ParticipantList({ open, onClose }: Props) {
  const participants = useRoomStore((s) => s.participants)

  if (!open) return null

  return (
    <>
      <div
        className="fixed inset-0 z-40 bg-black/50"
        onClick={onClose}
      />
      <div className="fixed inset-y-0 right-0 z-50 w-80 border-l border-zinc-800 bg-zinc-950 p-6 shadow-xl">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-2">
            <Users className="h-5 w-5 text-zinc-400" />
            <h2 className="text-lg font-semibold">Participants</h2>
            <span className="text-sm text-zinc-500">
              ({participants.length})
            </span>
          </div>
          <Button variant="ghost" size="sm" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </div>

        <div className="space-y-2">
          {participants.map((p) => (
            <div
              key={p.userId}
              className="flex items-center gap-3 rounded-lg bg-zinc-900 px-3 py-2"
            >
              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-zinc-800 text-sm font-medium">
                {p.displayName.charAt(0).toUpperCase()}
              </div>
              <div className="flex-1 min-w-0">
                <p className="truncate text-sm font-medium">{p.displayName}</p>
                <p className="text-xs text-zinc-500">
                  {p.connected ? 'Connected' : 'Disconnected'}
                </p>
              </div>
              <span
                className={`h-2 w-2 rounded-full ${
                  p.connected ? 'bg-emerald-500' : 'bg-zinc-600'
                }`}
              />
            </div>
          ))}
        </div>
      </div>
    </>
  )
}

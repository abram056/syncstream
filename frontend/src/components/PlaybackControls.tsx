import { Play, Pause, SkipBack, SkipForward, Users } from 'lucide-react'
import { Button } from './ui/button'

interface Props {
  isPlaying: boolean
  onTogglePlay: () => void
  onSeekBack: () => void
  onSeekForward: () => void
  currentTime: number
  duration: number
  onShowParticipants: () => void
}

function formatTime(seconds: number): string {
  if (!isFinite(seconds)) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

export function PlaybackControls({
  isPlaying,
  onTogglePlay,
  onSeekBack,
  onSeekForward,
  currentTime,
  duration,
  onShowParticipants,
}: Props) {
  return (
    <div className="flex items-center justify-center gap-4 px-4 py-3">
      <Button variant="ghost" size="sm" onClick={onSeekBack}>
        <SkipBack className="h-5 w-5" />
      </Button>

      <Button
        variant="secondary"
        size="md"
        className="h-12 w-12 rounded-full"
        onClick={onTogglePlay}
      >
        {isPlaying ? (
          <Pause className="h-5 w-5" />
        ) : (
          <Play className="h-5 w-5 ml-0.5" />
        )}
      </Button>

      <Button variant="ghost" size="sm" onClick={onSeekForward}>
        <SkipForward className="h-5 w-5" />
      </Button>

      <span className="text-xs text-zinc-500 tabular-nums min-w-[80px] text-right">
        {formatTime(currentTime)} / {formatTime(duration)}
      </span>

      <Button
        variant="ghost"
        size="sm"
        className="md:hidden"
        onClick={onShowParticipants}
      >
        <Users className="h-4 w-4" />
      </Button>
    </div>
  )
}

import { useRef, useState, useCallback } from 'react'
import { usePlaybackStore } from '../stores/playbackStore'
import { usePlaybackSync } from '../hooks/usePlaybackSync'
import { PlaybackControls } from './PlaybackControls'

interface Props {
  mediaUrl: string
  onTogglePlay: (position?: number) => void
  onSeek: (position: number) => void
  onShowParticipants: () => void
}

export function VideoPlayer({
  mediaUrl,
  onTogglePlay,
  onSeek,
  onShowParticipants,
}: Props) {
  const videoRef = useRef<HTMLVideoElement>(null)
  const [currentTime, setCurrentTime] = useState(0)
  const [duration, setDuration] = useState(0)
  const [seeking, setSeeking] = useState(false)
  const isPlaying = usePlaybackStore((s) => s.isPlaying)
  const lastSyncTime = usePlaybackStore((s) => s.lastSyncTime)

  usePlaybackSync(videoRef)

  const handleTogglePlay = useCallback(() => {
    if (videoRef.current) {
      const position = videoRef.current.currentTime
      onTogglePlay(position)
    }
  }, [onTogglePlay])

  const handleTimeUpdate = useCallback(() => {
    if (videoRef.current && !seeking) {
      setCurrentTime(videoRef.current.currentTime)
    }
  }, [seeking])

  const handleLoadedMetadata = useCallback(() => {
    if (videoRef.current) {
      setDuration(videoRef.current.duration)
    }
  }, [])

  const handleSeekBack = useCallback(() => {
    if (videoRef.current) {
      const newTime = Math.max(0, videoRef.current.currentTime - 10)
      videoRef.current.currentTime = newTime
      setCurrentTime(newTime)
      onSeek(newTime)
    }
  }, [onSeek])

  const handleSeekForward = useCallback(() => {
    if (videoRef.current) {
      const newTime = Math.min(
        videoRef.current.duration,
        videoRef.current.currentTime + 10
      )
      videoRef.current.currentTime = newTime
      setCurrentTime(newTime)
      onSeek(newTime)
    }
  }, [onSeek])

  // Show "Waiting for sync" overlay if not yet synced
  const showSyncOverlay = lastSyncTime === null

  return (
    <div className="relative flex flex-col">
      <div className="relative bg-black">
        <video
          ref={videoRef}
          src={mediaUrl}
          className="w-full max-h-[70vh] object-contain"
          onTimeUpdate={handleTimeUpdate}
          onLoadedMetadata={handleLoadedMetadata}
          onSeeking={() => setSeeking(true)}
          onSeeked={() => setSeeking(false)}
          playsInline
          preload="metadata"
        />

        {showSyncOverlay && (
          <div className="absolute inset-0 flex items-center justify-center bg-black/60">
            <div className="text-center space-y-2">
              <div className="h-8 w-8 animate-spin rounded-full border-2 border-blue-500 border-t-transparent mx-auto" />
              <p className="text-sm text-zinc-400">Synchronizing playback...</p>
            </div>
          </div>
        )}
      </div>

      <PlaybackControls
        isPlaying={isPlaying}
        onTogglePlay={handleTogglePlay}
        onSeekBack={handleSeekBack}
        onSeekForward={handleSeekForward}
        currentTime={currentTime}
        duration={duration}
        onShowParticipants={onShowParticipants}
      />
    </div>
  )
}

import { useEffect, useRef, useCallback } from 'react'
import { usePlaybackStore } from '../stores/playbackStore'

export function usePlaybackSync(videoRef: React.RefObject<HTMLVideoElement | null>) {
  const isPlaying = usePlaybackStore((s) => s.isPlaying)
  const position = usePlaybackStore((s) => s.position)
  const updatedAt = usePlaybackStore((s) => s.updatedAt)
  const lastSyncTime = usePlaybackStore((s) => s.lastSyncTime)
  const syncingRef = useRef(false)

  const effectivePosition = useCallback(() => {
    if (!updatedAt || !lastSyncTime) return position
    if (!isPlaying) return position
    const elapsed = (Date.now() - updatedAt) / 1000
    return elapsed > 0 ? position + elapsed : position
  }, [isPlaying, position, updatedAt, lastSyncTime])

  // When server says play/pause, follow
  useEffect(() => {
    const video = videoRef.current
    if (!video || !lastSyncTime) return

    const needsSeek = Math.abs(video.currentTime - effectivePosition()) > 1.0
    if (needsSeek) {
      syncingRef.current = true
      video.currentTime = effectivePosition()

      const onSeeked = () => {
        syncingRef.current = false
        video.removeEventListener('seeked', onSeeked)
      }
      video.addEventListener('seeked', onSeeked)
    }

    if (isPlaying) {
      video.play().catch(() => { })
    } else {
      video.pause()
    }
  }, [isPlaying, lastSyncTime, videoRef, effectivePosition])

  // Follow drifting position while playing
  useEffect(() => {
    if (!isPlaying || !lastSyncTime) return
    const interval = setInterval(() => {
      const video = videoRef.current
      if (!video || syncingRef.current) return
      const target = effectivePosition()
      const drift = Math.abs(video.currentTime - target)
      if (drift > 1) {
        video.currentTime = target
      }
    }, 5000)
    return () => clearInterval(interval)
  }, [isPlaying, lastSyncTime, videoRef, effectivePosition])

}

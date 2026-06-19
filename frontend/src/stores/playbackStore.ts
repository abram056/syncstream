import { create } from 'zustand'

interface PlaybackStore {
  isPlaying: boolean
  position: number
  updatedBy: string | null
  updatedAt: number | null
  lastSyncTime: number | null

  setPlaying: (playing: boolean) => void
  setPosition: (position: number) => void
  syncState: (data: {
    isPlaying: boolean
    position: number
    updatedBy: string
    updatedAt: number
  }) => void
  reset: () => void
}

export const usePlaybackStore = create<PlaybackStore>((set) => ({
  isPlaying: false,
  position: 0,
  updatedBy: null,
  updatedAt: null,
  lastSyncTime: null,

  setPlaying: (playing) => set({ isPlaying: playing }),
  setPosition: (position) => set({ position }),

  syncState: (data) =>
    set({
      isPlaying: data.isPlaying,
      position: data.position,
      updatedBy: data.updatedBy,
      updatedAt: data.updatedAt,
      lastSyncTime: Date.now(),
    }),

  reset: () =>
    set({
      isPlaying: false,
      position: 0,
      updatedBy: null,
      updatedAt: null,
      lastSyncTime: null,
    }),
}))

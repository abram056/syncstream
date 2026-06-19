import { create } from 'zustand'
import type { Participant, RoomView } from '../types'

interface RoomStore {
  roomId: string | null
  mediaUrl: string | null
  roomStatus: string | null
  view: RoomView
  participants: Participant[]
  localDisplayName: string

  setRoomId: (id: string) => void
  setMediaUrl: (url: string) => void
  setRoomStatus: (status: string) => void
  setView: (view: RoomView) => void
  setLocalDisplayName: (name: string) => void
  addParticipant: (userId: string, displayName: string) => void
  removeParticipant: (userId: string) => void
  setParticipantConnected: (userId: string, connected: boolean) => void
  reset: () => void
  resetRuntime: () => void
}

const initialState = {
  roomId: null,
  mediaUrl: null,
  roomStatus: null,
  view: 'lobby' as RoomView,
  participants: [] as Participant[],
  localDisplayName: '',
}

export const useRoomStore = create<RoomStore>((set) => ({
  ...initialState,

  setRoomId: (id) => set({ roomId: id }),
  setMediaUrl: (url) => set({ mediaUrl: url }),
  setRoomStatus: (status) => set({ roomStatus: status }),
  setView: (view) => set({ view }),
  setLocalDisplayName: (name) => set({ localDisplayName: name }),

  addParticipant: (userId, displayName) =>
    set((state) => {
      if (state.participants.some((p) => p.userId === userId)) return state
      return {
        participants: [
          ...state.participants,
          { userId, displayName, connected: true },
        ],
      }
    }),

  removeParticipant: (userId) =>
    set((state) => ({
      participants: state.participants.filter((p) => p.userId !== userId),
    })),

  setParticipantConnected: (userId, connected) =>
    set((state) => ({
      participants: state.participants.map((p) =>
        p.userId === userId ? { ...p, connected } : p
      ),
    })),

  reset: () => set(initialState),
  resetRuntime: () =>
    set({
      roomStatus: null,
      view: 'lobby' as RoomView,
      participants: [],
    }),
}))

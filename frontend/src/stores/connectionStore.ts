import { create } from 'zustand'
import type { ConnectionStatus } from '../services/websocket/manager'

interface ConnectionStore {
  status: ConnectionStatus
  reconnectCount: number
  setStatus: (status: ConnectionStatus) => void
  incrementReconnect: () => void
  resetReconnect: () => void
}

export const useConnectionStore = create<ConnectionStore>((set) => ({
  status: 'disconnected',
  reconnectCount: 0,

  setStatus: (status) => set({ status }),
  incrementReconnect: () =>
    set((s) => ({ reconnectCount: s.reconnectCount + 1 })),
  resetReconnect: () => set({ reconnectCount: 0 }),
}))

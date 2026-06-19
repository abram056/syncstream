import axios from 'axios'
import type {
  CreateRoomRequest,
  CreateRoomResponse,
  GetRoomResponse,
} from '../../types'

const http = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
})

export async function createRoom(
  data: CreateRoomRequest
): Promise<CreateRoomResponse> {
  const res = await http.post<CreateRoomResponse>('/rooms', data)
  return res.data
}

export async function getRoom(roomId: string): Promise<GetRoomResponse> {
  const res = await http.get<GetRoomResponse>(`/rooms/${roomId}`)
  return res.data
}

export async function healthCheck(): Promise<boolean> {
  try {
    await axios.get('/health')
    return true
  } catch {
    return false
  }
}

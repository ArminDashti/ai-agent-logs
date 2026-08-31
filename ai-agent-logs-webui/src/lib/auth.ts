export type User = {
  id: string
  email: string
  username: string
  created_at: string
  updated_at: string
}

export type AuthResponse = {
  token: string
  user: User
}

export type SessionRow = {
  id: string
  agent: string
  title: string
  models: string
  project: string
  app: string
  prompt_count: number
  rate: number
  device: string
  duration_ms: number
  token_total: number
  started_at: string
  ended_at: string
  created_at: string
}

export type SessionPromptRow = {
  id: string
  session_id: string
  prompt: string
  mode: string
  response: string
  rate: number
  duration_ms: number
  model?: string
  project?: string
  app?: string
  user_ip?: string
  token_usage?: string
  refined_prompt?: string
  is_accepted?: boolean
  date?: string
  time?: string
  device?: string
  session_title?: string
  created_at: string
}

export type SessionDetail = SessionRow & {
  prompts: SessionPromptRow[]
}

export type DashboardStats = {
  session_count: number
  avg_rate: number
  total_duration_ms: number
  token_total: number
}

const TOKEN_KEY = 'ai-agent-token'
const USER_KEY = 'ai-agent-user'

export const API_BASE = (() => {
  const raw = import.meta.env.VITE_API_BASE_URL as string | undefined
  if (raw !== undefined && raw !== '') return raw.replace(/\/$/, '')
  return (import.meta.env.BASE_URL || '/').replace(/\/$/, '')
})()

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function getStoredUser(): User | null {
  const raw = localStorage.getItem(USER_KEY)
  if (!raw) return null
  try {
    return JSON.parse(raw) as User
  } catch {
    return null
  }
}

export function setSession(token: string, user: User): void {
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(USER_KEY, JSON.stringify(user))
}

export function clearSession(): void {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
}

async function apiFetch<T>(path: string, options: RequestInit = {}, auth = false): Promise<T> {
  const headers = new Headers(options.headers)
  if (!headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }
  if (auth) {
    const token = getToken()
    if (!token) throw new Error('Not authenticated')
    headers.set('Authorization', `Bearer ${token}`)
  }

  const response = await fetch(`${API_BASE}${path}`, { ...options, headers })
  if (!response.ok) {
    let message = `Request failed (${response.status})`
    try {
      const data = (await response.json()) as { error?: string }
      if (data.error) message = data.error
    } catch {
      /* ignore */
    }
    throw new Error(message)
  }
  return response.json() as Promise<T>
}

export function login(body: { username: string; password: string }): Promise<AuthResponse> {
  return apiFetch<AuthResponse>('/api/v1/auth/login', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export function fetchSessions(): Promise<SessionRow[]> {
  return apiFetch<SessionRow[]>('/api/v1/sessions', {}, true)
}

export function fetchSession(id: string): Promise<SessionDetail> {
  return apiFetch<SessionDetail>(`/api/v1/sessions/${encodeURIComponent(id)}`, {}, true)
}

export function fetchDashboardStats(): Promise<DashboardStats> {
  return apiFetch<DashboardStats>('/api/v1/dashboard/stats', {}, true)
}

export function formatDuration(ms: number): string {
  if (ms < 1000) return `${ms}ms`
  const totalSec = Math.round(ms / 1000)
  if (totalSec < 60) return `${totalSec}s`
  const min = Math.floor(totalSec / 60)
  const sec = totalSec % 60
  return sec === 0 ? `${min}m` : `${min}m ${sec}s`
}

export function formatDate(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return '—'
  return d.toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: '2-digit',
  })
}

export function formatTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return '—'
  return d.toLocaleTimeString(undefined, {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

/** Format as YYYY-MM-DD HH:MM */
export function formatDateTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return '—'
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${day} ${hh}:${mm}`
}

export function formatRate(rate: number): string {
  if (!Number.isFinite(rate)) return '—'
  return Number.isInteger(rate) ? String(rate) : rate.toFixed(1)
}

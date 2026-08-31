export type PromptRow = {
  id: string
  project: string
  title: string
  mode: string
  duration_ms: number
  device: string
  rate: number
  occurred_at: string
  prompt: string
  created_at: string
}

export type SessionRow = {
  id: string
  agent: string
  title: string
  models: string
  project: string
  prompt_count: number
  rate: number
  device: string
  duration_ms: number
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
  model: string
  project: string
  user_ip: string
  token_usage: string
  refined_prompt: string
  is_accepted: boolean
  date: string
  time: string
  device?: string
  session_title?: string
  created_at: string
}

export type SessionDetail = SessionRow & {
  prompts: SessionPromptRow[]
}

export type CreatePromptInput = {
  project: string
  title: string
  mode: string
  duration_ms: number
  device: string
  rate: number
  prompt: string
  occurred_at?: string
}

export type LogAgentResponseInput = {
  model: string
  rate: number
  date: string
  time: string
  device: string
  user_ip?: string
  duration_ms: number
  project: string
  app?: string
  token_usage?: string
  user_prompt: string
  refined_prompt?: string
  is_accepted: boolean
  agent_response: string
  session?: string
  title_of_session: string
  mode?: string
  agent?: string
}

export type AgentTurnLog = {
  session: string
  title_of_session: string
  turn: SessionPromptRow
}

export type LoginResult = {
  token: string
  user: {
    id: string
    email: string
    username: string
    created_at: string
    updated_at: string
  }
}

function envOr(key: string, fallback: string): string {
  const v = process.env[key]
  return v && v.trim() !== '' ? v.trim() : fallback
}

export class AiAgentLogsClient {
  readonly baseUrl: string
  readonly username: string
  readonly password: string
  private token: string | null = null

  constructor(opts?: { baseUrl?: string; username?: string; password?: string }) {
    this.baseUrl = (opts?.baseUrl ?? envOr('AI_AGENT_LOGS_API_URL', 'http://127.0.0.1:8120')).replace(
      /\/$/,
      '',
    )
    this.username = opts?.username ?? envOr('AI_AGENT_LOGS_USERNAME', 'armin')
    this.password = opts?.password ?? envOr('AI_AGENT_LOGS_PASSWORD', 'dopadopa123')
  }

  private async request<T>(
    path: string,
    options: RequestInit = {},
    auth = false,
  ): Promise<T> {
    const headers = new Headers(options.headers)
    if (!headers.has('Content-Type') && options.body) {
      headers.set('Content-Type', 'application/json')
    }
    if (auth) {
      const token = await this.ensureToken()
      headers.set('Authorization', `Bearer ${token}`)
    }

    const response = await fetch(`${this.baseUrl}${path}`, { ...options, headers })
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
    if (response.status === 204) {
      return undefined as T
    }
    return (await response.json()) as T
  }

  async ensureToken(): Promise<string> {
    if (this.token) return this.token
    const result = await this.login()
    return result.token
  }

  async health(): Promise<{ status: string }> {
    return this.request<{ status: string }>('/health')
  }

  async login(username?: string, password?: string): Promise<LoginResult> {
    const result = await this.request<LoginResult>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({
        username: username ?? this.username,
        password: password ?? this.password,
      }),
    })
    this.token = result.token
    return result
  }

  async listPrompts(): Promise<PromptRow[]> {
    return this.request<PromptRow[]>('/api/v1/prompts', {}, true)
  }

  async createPrompt(input: CreatePromptInput): Promise<PromptRow> {
    return this.request<PromptRow>(
      '/api/v1/prompts',
      {
        method: 'POST',
        body: JSON.stringify(input),
      },
      true,
    )
  }

  async listSessions(): Promise<SessionRow[]> {
    return this.request<SessionRow[]>('/api/v1/sessions', {}, true)
  }

  async getSession(id: string): Promise<SessionDetail> {
    return this.request<SessionDetail>(
      `/api/v1/sessions/${encodeURIComponent(id)}`,
      {},
      true,
    )
  }

  async logAgentResponse(input: LogAgentResponseInput): Promise<AgentTurnLog> {
    return this.request<AgentTurnLog>(
      '/api/v1/agent-responses',
      {
        method: 'POST',
        body: JSON.stringify(input),
      },
      true,
    )
  }
}

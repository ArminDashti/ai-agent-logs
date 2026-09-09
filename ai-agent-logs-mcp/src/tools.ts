import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js'
import { z } from 'zod'
import { AiAgentLogsClient } from './client.js'

function textResult(data: unknown) {
  return {
    content: [{ type: 'text' as const, text: JSON.stringify(data, null, 2) }],
  }
}

function errorResult(err: unknown) {
  const message = err instanceof Error ? err.message : String(err)
  return {
    content: [{ type: 'text' as const, text: message }],
    isError: true as const,
  }
}

export function registerTools(server: McpServer, client: AiAgentLogsClient): void {
  server.tool(
    'health_check',
    'Check ai-agent-logs-api health (GET /health).',
    {},
    async () => {
      try {
        return textResult(await client.health())
      } catch (err) {
        return errorResult(err)
      }
    },
  )

  server.tool(
    'login',
    'Login to ai-agent-logs-api and cache the JWT. Optional username/password override env defaults.',
    {
      username: z.string().optional().describe('Username (default: AI_AGENT_LOGS_USERNAME or armin)'),
      password: z.string().optional().describe('Password (default: AI_AGENT_LOGS_PASSWORD)'),
    },
    async ({ username, password }) => {
      try {
        const result = await client.login(username, password)
        return textResult({
          ok: true,
          user: result.user,
          token_preview: `${result.token.slice(0, 12)}…`,
        })
      } catch (err) {
        return errorResult(err)
      }
    },
  )

  server.tool(
    'list_prompts',
    'List stored prompts newest first (GET /api/v1/prompts). Auto-logs in if needed.',
    {},
    async () => {
      try {
        return textResult(await client.listPrompts())
      } catch (err) {
        return errorResult(err)
      }
    },
  )

  server.tool(
    'create_prompt',
    'Create a prompt log row (POST /api/v1/prompts). Auto-logs in if needed.',
    {
      project: z.string().describe('Project name'),
      title: z.string().describe('Short title'),
      mode: z.string().describe('Agent mode (e.g. agent, ask, plan)'),
      duration_ms: z.number().int().min(0).describe('Duration in milliseconds'),
      device: z.string().describe('Device / hostname'),
      rate: z.number().int().min(0).max(10).describe('Integer prompt rate 0–10'),
      prompt: z.string().describe('Full prompt text'),
      occurred_at: z
        .string()
        .optional()
        .describe('Optional RFC3339 timestamp; defaults to now on the API'),
    },
    async (args) => {
      try {
        const body = {
          project: args.project,
          title: args.title,
          mode: args.mode,
          duration_ms: args.duration_ms,
          device: args.device,
          rate: args.rate,
          prompt: args.prompt,
          ...(args.occurred_at ? { occurred_at: args.occurred_at } : {}),
        }
        return textResult(await client.createPrompt(body))
      } catch (err) {
        return errorResult(err)
      }
    },
  )

  server.tool(
    'log_agent_response',
    'Log one full agent turn to ai-agent-logs-api (POST /api/v1/agent-responses). Creates or reuses a session (by required session id), stores the turn, and mirrors a prompts row. Auto-logs in if needed.',
    {
      model: z.string().describe('Model name (e.g. Composer, Auto)'),
      rate: z.number().int().min(0).max(10).describe('Integer prompt rate 0–10'),
      date: z.string().describe('Local date YYYY-MM-DD'),
      time: z.string().describe('Local time HH:MM:SS'),
      device: z.string().describe('Device / hostname'),
      user_ip: z.string().describe('User IP'),
      duration_ms: z.number().int().min(0).describe('Turn duration in milliseconds'),
      project: z.string().describe('Project / workspace name'),
      app: z.string().describe('Client app name (e.g. cursor)'),
      input_token: z.number().int().min(0).describe('Input / context token count'),
      output_token: z.number().int().min(0).describe('Output token count'),
      user_prompt: z.string().describe('Raw human prompt'),
      refined_prompt: z.string().describe('Transformed / refined prompt; empty string if none'),
      is_accepted: z.boolean().describe('Whether Human prompt is accepted'),
      agent_response: z.string().describe('Final agent user-facing response text'),
      session: z.string().describe('Session id (UUID); creates the session when unknown'),
      title_of_session: z.string().describe('Session title (chat title)'),
      mode: z.string().describe('Agent mode (e.g. agent, ask, plan)'),
      agent: z.string().describe('Agent product name (e.g. Cursor)'),
    },
    async (args) => {
      try {
        const body = {
          model: args.model,
          rate: args.rate,
          date: args.date,
          time: args.time,
          device: args.device,
          user_ip: args.user_ip,
          duration_ms: args.duration_ms,
          project: args.project,
          app: args.app,
          input_token: args.input_token,
          output_token: args.output_token,
          user_prompt: args.user_prompt,
          refined_prompt: args.refined_prompt,
          is_accepted: args.is_accepted,
          agent_response: args.agent_response,
          session: args.session,
          title_of_session: args.title_of_session,
          mode: args.mode,
          agent: args.agent,
        }
        return textResult(await client.logAgentResponse(body))
      } catch (err) {
        return errorResult(err)
      }
    },
  )

  server.tool(
    'list_sessions',
    'List agent conversation sessions newest first (GET /api/v1/sessions). Auto-logs in if needed.',
    {},
    async () => {
      try {
        return textResult(await client.listSessions())
      } catch (err) {
        return errorResult(err)
      }
    },
  )

  server.tool(
    'get_session',
    'Get one session and its prompt turns (GET /api/v1/sessions/:id). Auto-logs in if needed.',
    {
      id: z.string().describe('Session id'),
    },
    async ({ id }) => {
      try {
        return textResult(await client.getSession(id))
      } catch (err) {
        return errorResult(err)
      }
    },
  )
}

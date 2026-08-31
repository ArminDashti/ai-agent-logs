# ai-agent-logs-mcp

MCP server so Cursor agents can talk to [ai-agent-logs-api](https://github.com/ArminDashti/ai-agent-logs-api) (same data shown in [ai-agent-logs-webui](https://github.com/ArminDashti/ai-agent-logs-webui)).

## Prerequisites

1. Run the API (local `http://127.0.0.1:8120` or published `https://ai-agent-log-api.xaigrok.ir`).
2. Node.js 18+.

## Setup

```powershell
cd C:\Users\armin\GitHub\ai-agent-logs-mcp
copy .env.example .env
npm install
npm run build
```

## Cursor (`~/.cursor/mcp.json`)

```json
"ai-agent-logs-mcp": {
  "command": "node",
  "args": ["C:/Users/armin/GitHub/ai-agent-logs-mcp/dist/index.js"],
  "env": {
    "AI_AGENT_LOGS_API_URL": "https://ai-agent-log-api.xaigrok.ir",
    "AI_AGENT_LOGS_USERNAME": "armin",
    "AI_AGENT_LOGS_PASSWORD": "dopadopa123"
  }
}
```

Restart MCP / reload Cursor after changing `mcp.json`.

## Tools

| Tool | API |
|------|-----|
| `health_check` | `GET /health` |
| `login` | `POST /api/v1/auth/login` |
| `list_prompts` | `GET /api/v1/prompts` |
| `create_prompt` | `POST /api/v1/prompts` |
| `log_agent_response` | `POST /api/v1/agent-responses` |
| `list_sessions` | `GET /api/v1/sessions` |
| `get_session` | `GET /api/v1/sessions/:id` |

### `log_agent_response`

Full turn fields sent to the API: `model`, `rate`, `date`, `time`, `device`, `user_ip`, `duration_ms`, `project`, `app` (default `cursor`), `token_usage`, `user_prompt`, `refined_prompt`, `is_accepted`, `agent_response`, `session`, `title_of_session`.

Authenticated tools auto-login with env credentials when no JWT is cached.

## Env

| Variable | Default |
|----------|---------|
| `AI_AGENT_LOGS_API_URL` | `http://127.0.0.1:8120` |
| `AI_AGENT_LOGS_USERNAME` | `armin` |
| `AI_AGENT_LOGS_PASSWORD` | `dopadopa123` |

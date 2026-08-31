# ai-agent-api

Gin + PostgreSQL API for AI Agent prompt history.

## Local setup

```powershell
docker compose -f docker-compose.local.yml up -d
copy .env.example .env
go mod tidy
go run ./cmd/server
```

- API: `http://localhost:8120`
- Postgres: `localhost:5434` (user/pass/db `ai_agent`)
- Default login: `armin` / `dopadopa123`

## Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | no | Health check |
| POST | `/api/v1/auth/login` | no | JWT login |
| GET | `/api/v1/prompts` | yes | List prompts |
| POST | `/api/v1/prompts` | yes | Create prompt |

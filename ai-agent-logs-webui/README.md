# ai-agent-webui

Vue 3 + Vite + Tailwind + shadcn-style UI for AI Agent prompt history.

## Local setup

```powershell
npm install
npm run dev
```

- Dev server: `http://localhost:5175` (proxies `/api` to `http://127.0.0.1:8120`)
- Font: Inter via `@fontsource/inter`
- Default login: `armin` / `dopadopa123`

## Pages

- `/prompts` — grid: Project, Title, Mode, Duration, Device, Rate, Date, Time, Prompt
- `/login` — JWT login
- `/about` — About Me

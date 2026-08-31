package handlers

import (
	"database/sql"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ArminDashti/ai-agent-api/internal/models"
	"github.com/gin-gonic/gin"
)

const sessionSelectSQL = `
		SELECT
			s.id,
			s.agent,
			s.title,
			s.models,
			COALESCE(s.project, '') AS project,
			COALESCE(s.app, 'cursor') AS app,
			COALESCE(COUNT(sp.id), 0)::int AS prompt_count,
			COALESCE(AVG(sp.rate), 0)::float8 AS rate,
			s.device,
			COALESCE(SUM(sp.duration_ms), 0)::int AS duration_ms,
			COALESCE(SUM(parse_token_usage(sp.token_usage)), 0)::int AS token_total,
			COALESCE(MIN(sp.created_at), s.created_at) AS started_at,
			COALESCE(MAX(sp.created_at), s.created_at) AS ended_at,
			s.created_at
		FROM sessions s
		LEFT JOIN session_prompts sp ON sp.session_id = s.id
`

// ListSessions returns agent conversation sessions newest first.
func (h *Handler) ListSessions(c *gin.Context) {
	rows, err := h.db.QueryContext(c.Request.Context(), sessionSelectSQL+`
		GROUP BY s.id
		ORDER BY s.created_at DESC
	`)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not list sessions")
		return
	}
	defer rows.Close()

	out := make([]models.Session, 0)
	for rows.Next() {
		var s models.Session
		if err := rows.Scan(
			&s.ID, &s.Agent, &s.Title, &s.Models, &s.Project, &s.App, &s.PromptCount,
			&s.Rate, &s.Device, &s.DurationMs, &s.TokenTotal, &s.StartedAt, &s.EndedAt, &s.CreatedAt,
		); err != nil {
			writeError(c, http.StatusInternalServerError, "could not read sessions")
			return
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		writeError(c, http.StatusInternalServerError, "could not read sessions")
		return
	}
	c.JSON(http.StatusOK, out)
}

// GetSession returns one session and its prompt turns.
func (h *Handler) GetSession(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		writeError(c, http.StatusBadRequest, "session id is required")
		return
	}

	var detail models.SessionDetail
	err := h.db.QueryRowContext(c.Request.Context(), sessionSelectSQL+`
		WHERE s.id = $1
		GROUP BY s.id
	`, id).Scan(
		&detail.ID, &detail.Agent, &detail.Title, &detail.Models, &detail.Project, &detail.App, &detail.PromptCount,
		&detail.Rate, &detail.Device, &detail.DurationMs, &detail.TokenTotal, &detail.StartedAt, &detail.EndedAt, &detail.CreatedAt,
	)
	if err == sql.ErrNoRows {
		writeError(c, http.StatusNotFound, "session not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load session")
		return
	}

	promptRows, err := h.db.QueryContext(c.Request.Context(), `
		SELECT
			id, session_id, prompt, mode, response, rate, duration_ms,
			COALESCE(model, ''), COALESCE(project, ''), COALESCE(app, 'cursor'), COALESCE(user_ip, ''),
			COALESCE(token_usage, ''), COALESCE(refined_prompt, ''), COALESCE(is_accepted, TRUE),
			COALESCE(to_char(log_date, 'YYYY-MM-DD'), ''),
			COALESCE(to_char(log_time, 'HH24:MI:SS'), ''),
			created_at
		FROM session_prompts
		WHERE session_id = $1
		ORDER BY created_at ASC
	`, id)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not list session prompts")
		return
	}
	defer promptRows.Close()

	detail.Prompts = make([]models.SessionPrompt, 0)
	for promptRows.Next() {
		var p models.SessionPrompt
		if err := promptRows.Scan(
			&p.ID, &p.SessionID, &p.Prompt, &p.Mode, &p.Response,
			&p.Rate, &p.DurationMs, &p.Model, &p.Project, &p.App, &p.UserIP,
			&p.TokenUsage, &p.RefinedPrompt, &p.IsAccepted, &p.LogDate, &p.LogTime,
			&p.CreatedAt,
		); err != nil {
			writeError(c, http.StatusInternalServerError, "could not read session prompts")
			return
		}
		detail.Prompts = append(detail.Prompts, p)
	}
	if err := promptRows.Err(); err != nil {
		writeError(c, http.StatusInternalServerError, "could not read session prompts")
		return
	}

	c.JSON(http.StatusOK, detail)
}

// DashboardStats returns KPI aggregates for the dashboard page.
func (h *Handler) DashboardStats(c *gin.Context) {
	var stats models.DashboardStats
	err := h.db.QueryRowContext(c.Request.Context(), `
		SELECT
			(SELECT COUNT(*)::int FROM sessions) AS session_count,
			COALESCE((SELECT AVG(rate)::float8 FROM session_prompts), 0) AS avg_rate,
			COALESCE((SELECT SUM(duration_ms)::int FROM session_prompts), 0) AS total_duration_ms,
			COALESCE((SELECT SUM(parse_token_usage(token_usage))::int FROM session_prompts), 0) AS token_total
	`).Scan(&stats.SessionCount, &stats.AvgRate, &stats.TotalDurationMs, &stats.TokenTotal)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load dashboard stats")
		return
	}
	c.JSON(http.StatusOK, stats)
}

type logAgentResponseRequest struct {
	Model          string `json:"model" binding:"required"`
	Rate           *int   `json:"rate" binding:"required"`
	Date           string `json:"date" binding:"required"`
	Time           string `json:"time" binding:"required"`
	Device         string `json:"device" binding:"required"`
	UserIP         string `json:"user_ip"`
	DurationMs     *int   `json:"duration_ms" binding:"required"`
	Project        string `json:"project" binding:"required"`
	App            string `json:"app"`
	TokenUsage     string `json:"token_usage"`
	UserPrompt     string `json:"user_prompt" binding:"required"`
	RefinedPrompt  string `json:"refined_prompt"`
	IsAccepted     *bool  `json:"is_accepted" binding:"required"`
	AgentResponse  string `json:"agent_response" binding:"required"`
	Session        string `json:"session"`
	TitleOfSession string `json:"title_of_session" binding:"required"`
	Mode           string `json:"mode"`
	Agent          string `json:"agent"`
}

// LogAgentResponse creates or reuses a session and stores one full agent turn.
func (h *Handler) LogAgentResponse(c *gin.Context) {
	var req logAgentResponseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}

	model := strings.TrimSpace(req.Model)
	device := strings.TrimSpace(req.Device)
	project := strings.TrimSpace(req.Project)
	userPrompt := strings.TrimSpace(req.UserPrompt)
	agentResponse := strings.TrimSpace(req.AgentResponse)
	sessionTitle := strings.TrimSpace(req.TitleOfSession)
	dateStr := strings.TrimSpace(req.Date)
	timeStr := strings.TrimSpace(req.Time)
	if model == "" || device == "" || project == "" || userPrompt == "" || agentResponse == "" || sessionTitle == "" {
		writeError(c, http.StatusBadRequest, "model, device, project, user_prompt, agent_response, and title_of_session are required")
		return
	}
	if req.DurationMs == nil || *req.DurationMs < 0 {
		writeError(c, http.StatusBadRequest, "duration_ms must be >= 0")
		return
	}
	if req.Rate == nil || *req.Rate < 0 || *req.Rate > 10 {
		writeError(c, http.StatusBadRequest, "rate must be between 0 and 10")
		return
	}
	if req.IsAccepted == nil {
		writeError(c, http.StatusBadRequest, "is_accepted is required")
		return
	}

	logDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(c, http.StatusBadRequest, "date must be YYYY-MM-DD")
		return
	}
	logTime, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		writeError(c, http.StatusBadRequest, "time must be HH:MM:SS")
		return
	}

	userIP := strings.TrimSpace(req.UserIP)
	if userIP == "" {
		userIP = clientIP(c)
	}
	tokenUsage := strings.TrimSpace(req.TokenUsage)
	refined := strings.TrimSpace(req.RefinedPrompt)
	mode := strings.TrimSpace(req.Mode)
	if mode == "" {
		mode = "agent"
	}
	agentName := strings.TrimSpace(req.Agent)
	if agentName == "" {
		agentName = "Cursor"
	}
	appName := strings.TrimSpace(req.App)
	if appName == "" {
		appName = "cursor"
	}

	tx, err := h.db.BeginTx(c.Request.Context(), nil)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not start transaction")
		return
	}
	defer func() { _ = tx.Rollback() }()

	sessionID := strings.TrimSpace(req.Session)
	var existingTitle string
	if sessionID != "" {
		err = tx.QueryRowContext(c.Request.Context(), `
			SELECT title FROM sessions WHERE id = $1
		`, sessionID).Scan(&existingTitle)
		if err == sql.ErrNoRows {
			writeError(c, http.StatusNotFound, "session not found")
			return
		}
		if err != nil {
			writeError(c, http.StatusInternalServerError, "could not load session")
			return
		}
		sessionTitle = existingTitle
		_, err = tx.ExecContext(c.Request.Context(), `
			UPDATE sessions
			SET models = CASE
				WHEN models = '' THEN $2
				WHEN position($2 in models) > 0 THEN models
				ELSE models || ', ' || $2
			END,
			project = CASE WHEN project = '' THEN $3 ELSE project END,
			app = CASE WHEN app = '' OR app = 'cursor' THEN $4 ELSE app END
			WHERE id = $1
		`, sessionID, model, project, appName)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "could not update session")
			return
		}
	} else {
		err = tx.QueryRowContext(c.Request.Context(), `
			INSERT INTO sessions (agent, title, models, device, project, app)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, title
		`, agentName, sessionTitle, model, device, project, appName).Scan(&sessionID, &sessionTitle)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "could not create session")
			return
		}
	}

	var turn models.SessionPrompt
	err = tx.QueryRowContext(c.Request.Context(), `
		INSERT INTO session_prompts (
			session_id, prompt, mode, response, rate, duration_ms,
			model, project, app, user_ip, token_usage, refined_prompt, is_accepted,
			log_date, log_time
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13,
			$14, $15
		)
		RETURNING
			id, session_id, prompt, mode, response, rate, duration_ms,
			COALESCE(model, ''), COALESCE(project, ''), COALESCE(app, 'cursor'), COALESCE(user_ip, ''),
			COALESCE(token_usage, ''), COALESCE(refined_prompt, ''), COALESCE(is_accepted, TRUE),
			COALESCE(to_char(log_date, 'YYYY-MM-DD'), ''),
			COALESCE(to_char(log_time, 'HH24:MI:SS'), ''),
			created_at
	`, sessionID, userPrompt, mode, agentResponse, *req.Rate, *req.DurationMs,
		model, project, appName, userIP, tokenUsage, refined, *req.IsAccepted,
		logDate, logTime,
	).Scan(
		&turn.ID, &turn.SessionID, &turn.Prompt, &turn.Mode, &turn.Response,
		&turn.Rate, &turn.DurationMs, &turn.Model, &turn.Project, &turn.App, &turn.UserIP,
		&turn.TokenUsage, &turn.RefinedPrompt, &turn.IsAccepted, &turn.LogDate, &turn.LogTime,
		&turn.CreatedAt,
	)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not save agent response")
		return
	}

	occurredAt := time.Date(logDate.Year(), logDate.Month(), logDate.Day(),
		logTime.Hour(), logTime.Minute(), logTime.Second(), 0, time.Local).UTC()
	_, err = tx.ExecContext(c.Request.Context(), `
		INSERT INTO prompts (project, title, mode, duration_ms, device, rate, occurred_at, prompt)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, project, sessionTitle, mode, *req.DurationMs, device, *req.Rate, occurredAt, userPrompt)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not mirror prompt row")
		return
	}

	if err := tx.Commit(); err != nil {
		writeError(c, http.StatusInternalServerError, "could not commit agent response")
		return
	}

	turn.Device = device
	turn.SessionTitle = sessionTitle
	c.JSON(http.StatusCreated, models.AgentTurnLog{
		SessionID:    sessionID,
		SessionTitle: sessionTitle,
		Turn:         turn,
	})
}

func clientIP(c *gin.Context) string {
	ip := strings.TrimSpace(c.ClientIP())
	if ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	return strings.TrimSpace(c.Request.RemoteAddr)
}

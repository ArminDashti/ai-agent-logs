package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/ArminDashti/ai-agent-api/internal/auth"
	"github.com/ArminDashti/ai-agent-api/internal/config"
	"github.com/ArminDashti/ai-agent-api/internal/models"
	"github.com/gin-gonic/gin"
)

// Handler holds shared dependencies.
type Handler struct {
	db  *sql.DB
	cfg config.Config
}

// New creates a Handler.
func New(db *sql.DB, cfg config.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

func writeError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

// Health returns a simple health payload.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates by username and returns a JWT.
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}
	username := strings.TrimSpace(req.Username)
	if username == "" {
		writeError(c, http.StatusBadRequest, "username is required")
		return
	}

	var user models.User
	err := h.db.QueryRowContext(c.Request.Context(), `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users WHERE username = $1
	`, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		writeError(c, http.StatusUnauthorized, "invalid username or password")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load user")
		return
	}
	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		writeError(c, http.StatusUnauthorized, "invalid username or password")
		return
	}

	token, err := auth.IssueToken(h.cfg.JWTSecret, user.ID, user.Username)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not issue token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"username":   user.Username,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

type createPromptRequest struct {
	Project    string  `json:"project" binding:"required"`
	Title      string  `json:"title" binding:"required"`
	Mode       string  `json:"mode" binding:"required"`
	DurationMs *int    `json:"duration_ms" binding:"required"`
	Device     string  `json:"device" binding:"required"`
	Rate       *int    `json:"rate" binding:"required"`
	OccurredAt *string `json:"occurred_at"`
	Prompt     string  `json:"prompt" binding:"required"`
}

// ListPrompts returns prompts newest first.
func (h *Handler) ListPrompts(c *gin.Context) {
	rows, err := h.db.QueryContext(c.Request.Context(), `
		SELECT id, project, title, mode, duration_ms, device, rate, occurred_at, prompt, created_at
		FROM prompts
		ORDER BY occurred_at DESC
	`)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not list prompts")
		return
	}
	defer rows.Close()

	out := make([]models.Prompt, 0)
	for rows.Next() {
		var p models.Prompt
		if err := rows.Scan(
			&p.ID, &p.Project, &p.Title, &p.Mode, &p.DurationMs, &p.Device,
			&p.Rate, &p.OccurredAt, &p.Prompt, &p.CreatedAt,
		); err != nil {
			writeError(c, http.StatusInternalServerError, "could not read prompts")
			return
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		writeError(c, http.StatusInternalServerError, "could not read prompts")
		return
	}
	c.JSON(http.StatusOK, out)
}

// CreatePrompt persists a prompt row.
func (h *Handler) CreatePrompt(c *gin.Context) {
	var req createPromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}

	project := strings.TrimSpace(req.Project)
	title := strings.TrimSpace(req.Title)
	mode := strings.TrimSpace(req.Mode)
	device := strings.TrimSpace(req.Device)
	prompt := strings.TrimSpace(req.Prompt)
	if project == "" || title == "" || mode == "" || device == "" || prompt == "" {
		writeError(c, http.StatusBadRequest, "project, title, mode, device, and prompt are required")
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

	occurredAt := time.Now().UTC()
	if req.OccurredAt != nil && strings.TrimSpace(*req.OccurredAt) != "" {
		parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.OccurredAt))
		if err != nil {
			writeError(c, http.StatusBadRequest, "occurred_at must be RFC3339")
			return
		}
		occurredAt = parsed
	}

	var p models.Prompt
	err := h.db.QueryRowContext(c.Request.Context(), `
		INSERT INTO prompts (project, title, mode, duration_ms, device, rate, occurred_at, prompt)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, project, title, mode, duration_ms, device, rate, occurred_at, prompt, created_at
	`, project, title, mode, *req.DurationMs, device, *req.Rate, occurredAt, prompt).Scan(
		&p.ID, &p.Project, &p.Title, &p.Mode, &p.DurationMs, &p.Device,
		&p.Rate, &p.OccurredAt, &p.Prompt, &p.CreatedAt,
	)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not save prompt")
		return
	}
	c.JSON(http.StatusCreated, p)
}

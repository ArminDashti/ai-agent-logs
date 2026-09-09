package models

import "time"

// User is an authenticated account.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Prompt is a stored AI Agent prompt row.
type Prompt struct {
	ID         string    `json:"id"`
	Project    string    `json:"project"`
	Title      string    `json:"title"`
	Mode       string    `json:"mode"`
	DurationMs int       `json:"duration_ms"`
	Device     string    `json:"device"`
	Rate       int       `json:"rate"`
	OccurredAt time.Time `json:"occurred_at"`
	Prompt     string    `json:"prompt"`
	CreatedAt  time.Time `json:"created_at"`
}

// Session is an agent conversation session summary row.
type Session struct {
	ID          string    `json:"id"`
	Agent       string    `json:"agent"`
	Title       string    `json:"title"`
	Models      string    `json:"models"`
	Project     string    `json:"project"`
	App         string    `json:"app"`
	PromptCount int       `json:"prompt_count"`
	Rate        float64   `json:"rate"`
	Device      string    `json:"device"`
	DurationMs  int       `json:"duration_ms"`
	TokenTotal  int       `json:"token_total"`
	StartedAt   time.Time `json:"started_at"`
	EndedAt     time.Time `json:"ended_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// SessionPrompt is one prompt turn inside a session.
type SessionPrompt struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id"`
	Prompt        string    `json:"prompt"`
	Mode          string    `json:"mode"`
	Response      string    `json:"response"`
	Rate          int       `json:"rate"`
	DurationMs    int       `json:"duration_ms"`
	Model         string    `json:"model"`
	Project       string    `json:"project"`
	App           string    `json:"app"`
	UserIP        string    `json:"user_ip"`
	InputToken    int       `json:"input_token"`
	OutputToken   int       `json:"output_token"`
	RefinedPrompt string    `json:"refined_prompt"`
	IsAccepted    bool      `json:"is_accepted"`
	LogDate       string    `json:"date"`
	LogTime       string    `json:"time"`
	Device        string    `json:"device,omitempty"`
	SessionTitle  string    `json:"session_title,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// AgentTurnLog is the payload returned after logging one agent response.
type AgentTurnLog struct {
	SessionID    string        `json:"session"`
	SessionTitle string        `json:"title_of_session"`
	Turn         SessionPrompt `json:"turn"`
}

// SessionDetail is a session with its prompt turns.
type SessionDetail struct {
	Session
	Prompts []SessionPrompt `json:"prompts"`
}

// DashboardStats is aggregated KPI data for the dashboard page.
type DashboardStats struct {
	SessionCount     int     `json:"session_count"`
	AvgRate          float64 `json:"avg_rate"`
	TotalDurationMs  int     `json:"total_duration_ms"`
	TokenTotal       int     `json:"token_total"`
}

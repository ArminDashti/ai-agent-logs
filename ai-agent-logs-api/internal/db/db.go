package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Open connects to PostgreSQL.
func Open(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

// Migrate runs SQL files in migrationsDir in lexical order.
func Migrate(db *sql.DB, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, name := range files {
		body, err := os.ReadFile(filepath.Join(migrationsDir, name))
		if err != nil {
			return err
		}
		if _, err := db.ExecContext(ctx, string(body)); err != nil {
			return fmt.Errorf("migration %s: %w", name, err)
		}
	}
	return nil
}

// Default local/demo login account.
const (
	DefaultUsername = "armin"
	DefaultPassword = "dopadopa123"
	DefaultEmail    = "armin@local"
)

// SeedDefaultUser upserts the default local login account.
func SeedDefaultUser(ctx context.Context, db *sql.DB, passwordHash string) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO users (email, password_hash, username)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE SET
			password_hash = EXCLUDED.password_hash,
			username = EXCLUDED.username,
			updated_at = NOW()
	`, DefaultEmail, passwordHash, DefaultUsername)
	return err
}

// SeedSamplePrompts inserts demo rows when the prompts table is empty.
func SeedSamplePrompts(ctx context.Context, db *sql.DB) error {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM prompts`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	_, err := db.ExecContext(ctx, `
		INSERT INTO prompts (project, title, mode, duration_ms, device, rate, occurred_at, prompt)
		VALUES
			('ai-agent', 'Scaffold twin repos', 'agent', 125000, 'ARMIN-DESKTOP', 8, NOW() - INTERVAL '2 hours',
			 'Create public GitHub repos ai-agent-api and ai-agent-webui with Gin/Postgres and Vue/Shadcn.'),
			('armin-command-center', 'Update projects inventory', 'ask', 42000, 'ARMIN-DESKTOP', 9, NOW() - INTERVAL '1 hour',
			 'Add ai-agent-api and ai-agent-webui rows to projects.md.'),
			('prompt-corrector', 'Review prompts grid', 'plan', 89000, 'ARMIN-DESKTOP', 7, NOW() - INTERVAL '30 minutes',
			 'Plan a prompts page with Project, Title, Mode, Duration, Device, Rate, Date, Time, Prompt.')
	`)
	return err
}

// SeedSampleSessions inserts demo session rows when the sessions table is empty.
func SeedSampleSessions(ctx context.Context, db *sql.DB) error {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sessions`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var cursorID, devinID string
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO sessions (agent, title, models, device, created_at)
		VALUES
			('Cursor', 'Rename cursor-ai to ai-agent', 'claude-opus-4, composer-2', 'ARMIN-DESKTOP', NOW() - INTERVAL '3 hours')
		RETURNING id
	`).Scan(&cursorID); err != nil {
		return err
	}
	if err := tx.QueryRowContext(ctx, `
		INSERT INTO sessions (agent, title, models, device, created_at)
		VALUES
			('Devin', 'Sessions page planning', 'gpt-5.2', 'MacBook-Pro', NOW() - INTERVAL '1 hour')
		RETURNING id
	`).Scan(&devinID); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO session_prompts (session_id, prompt, mode, response, rate, duration_ms, created_at)
		VALUES
			($1, 'Rename these repos to ai-agent-api and ai-agent-webui', 'agent',
			 'Renamed folders, module paths, and branding to AI Agent.', 8, 180000, NOW() - INTERVAL '3 hours'),
			($1, 'Add a sessions page with master and detail grids', 'agent',
			 'Designed session list columns and detail prompt turns.', 9, 95000, NOW() - INTERVAL '2 hours 40 minutes'),
			($2, 'Outline agent conversation logging schema', 'plan',
			 'Proposed sessions + session_prompts with aggregates.', 7, 61000, NOW() - INTERVAL '55 minutes'),
			($2, 'What columns belong on the sessions grid?', 'ask',
			 'Agent, Title, Models, Prompt count, Rate, Device, Duration, Created at.', 8, 24000, NOW() - INTERVAL '40 minutes')
	`, cursorID, devinID); err != nil {
		return err
	}

	return tx.Commit()
}

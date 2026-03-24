package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magomedcoder/gen/internal/domain"
)

type chatSessionSettingsRepository struct {
	db *pgxpool.Pool
}

func NewChatSessionSettingsRepository(db *pgxpool.Pool) domain.ChatSessionSettingsRepository {
	return &chatSessionSettingsRepository{db: db}
}

func (r *chatSessionSettingsRepository) GetBySessionID(ctx context.Context, sessionID int64) (*domain.ChatSessionSettings, error) {
	settings := &domain.ChatSessionSettings{
		SessionID: sessionID,
	}
	err := r.db.QueryRow(ctx, `
		SELECT system_prompt, stop_sequences, timeout_seconds, temperature, max_tokens, top_k, top_p, json_mode, json_schema, tools_json, profile
		FROM chat_session_settings
		WHERE session_id = $1
	`, sessionID).Scan(
		&settings.SystemPrompt,
		&settings.StopSequences,
		&settings.TimeoutSeconds,
		&settings.Temperature,
		&settings.MaxTokens,
		&settings.TopK,
		&settings.TopP,
		&settings.JSONMode,
		&settings.JSONSchema,
		&settings.ToolsJSON,
		&settings.Profile,
	)
	if err != nil {
		return settings, nil
	}

	return settings, nil
}

func (r *chatSessionSettingsRepository) Upsert(ctx context.Context, settings *domain.ChatSessionSettings) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO chat_session_settings (
		    session_id,
		    system_prompt,
		    stop_sequences,
		    timeout_seconds,
		    temperature,
		    max_tokens,
		    top_k,
		    top_p,
		    json_mode,
		    json_schema,
		    tools_json,
		    profile,
		    updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW())
		ON CONFLICT (session_id) DO UPDATE SET
		    system_prompt = EXCLUDED.system_prompt,
		    stop_sequences = EXCLUDED.stop_sequences,
		    timeout_seconds = EXCLUDED.timeout_seconds,
		    temperature = EXCLUDED.temperature,
		    max_tokens = EXCLUDED.max_tokens,
		    top_k = EXCLUDED.top_k,
		    top_p = EXCLUDED.top_p,
		    json_mode = EXCLUDED.json_mode,
		    json_schema = EXCLUDED.json_schema,
		    tools_json = EXCLUDED.tools_json,
		    profile = EXCLUDED.profile,
		    updated_at = NOW()
	`,
		settings.SessionID,
		settings.SystemPrompt,
		settings.StopSequences,
		settings.TimeoutSeconds,
		settings.Temperature,
		settings.MaxTokens,
		settings.TopK,
		settings.TopP,
		settings.JSONMode,
		settings.JSONSchema,
		settings.ToolsJSON,
		settings.Profile,
	)

	return err
}

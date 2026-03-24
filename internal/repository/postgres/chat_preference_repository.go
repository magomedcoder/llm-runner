package postgres

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magomedcoder/gen/internal/domain"
)

type chatPreferenceRepository struct {
	db *pgxpool.Pool
}

func NewChatPreferenceRepository(db *pgxpool.Pool) domain.ChatPreferenceRepository {
	return &chatPreferenceRepository{db: db}
}

func (r *chatPreferenceRepository) GetSelectedRunner(ctx context.Context, userID int) (string, error) {
	var selected string
	err := r.db.QueryRow(ctx, `
		SELECT selected_runner
		FROM user_chat_preferences
		WHERE user_id = $1
	`, userID).Scan(&selected)
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(selected), nil
}

func (r *chatPreferenceRepository) SetSelectedRunner(ctx context.Context, userID int, runner string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_chat_preferences (user_id, selected_runner, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id) DO UPDATE
		SET selected_runner = EXCLUDED.selected_runner,
		    updated_at = NOW()
	`, userID, strings.TrimSpace(runner))

	return err
}

func (r *chatPreferenceRepository) GetDefaultRunnerModel(ctx context.Context, userID int, runner string) (string, error) {
	var model string
	err := r.db.QueryRow(ctx, `
		SELECT model
		FROM user_runner_models
		WHERE user_id = $1 AND runner_address = $2
	`, userID, strings.TrimSpace(runner)).Scan(&model)
	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(model), nil
}

func (r *chatPreferenceRepository) SetDefaultRunnerModel(ctx context.Context, userID int, runner string, model string) error {
	runner = strings.TrimSpace(runner)
	model = strings.TrimSpace(model)
	if runner == "" {
		return nil
	}
	if model == "" {
		_, err := r.db.Exec(ctx, `
			DELETE FROM user_runner_models
			WHERE user_id = $1 AND runner_address = $2
		`, userID, runner)
		return err
	}

	_, err := r.db.Exec(ctx, `
		INSERT INTO user_runner_models (user_id, runner_address, model, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, runner_address) DO UPDATE
		SET model = EXCLUDED.model,
		    updated_at = NOW()
	`, userID, runner, model)

	return err
}

package postgres

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magomedcoder/gen/internal/domain"
)

type editorHistoryRepository struct {
	db *pgxpool.Pool
}

func NewEditorHistoryRepository(db *pgxpool.Pool) domain.EditorHistoryRepository {
	return &editorHistoryRepository{db: db}
}

func (r *editorHistoryRepository) Save(ctx context.Context, userID int, runner string, text string) error {
	if strings.TrimSpace(text) == "" {
		return nil
	}

	_, err := r.db.Exec(ctx, `
		INSERT INTO editor_text_history (user_id, runner, text, created_at)
		VALUES ($1, $2, $3, NOW())
	`, userID, strings.TrimSpace(runner), text)

	return err
}

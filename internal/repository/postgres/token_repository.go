package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/magomedcoder/gen/internal/domain"
)

type userSessionRepository struct {
	db *pgxpool.Pool
}

func NewUserSessionRepository(db *pgxpool.Pool) domain.TokenRepository {
	return &userSessionRepository{db: db}
}

func (u *userSessionRepository) Create(ctx context.Context, token *domain.Token) error {
	err := u.db.QueryRow(ctx,
		`
		INSERT INTO user_sessions (user_id, token, type, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`,
		token.UserId,
		token.Token,
		token.Type,
		token.ExpiresAt,
		token.CreatedAt,
	).Scan(&token.Id)

	return err
}

func (u *userSessionRepository) GetByToken(ctx context.Context, token string) (*domain.Token, error) {
	var t domain.Token
	err := u.db.QueryRow(ctx,
		`
		SELECT id, user_id, token, type, expires_at, created_at, deleted_at
		FROM user_sessions
		WHERE token = $1 AND deleted_at IS NULL
	`, token).Scan(
		&t.Id,
		&t.UserId,
		&t.Token,
		&t.Type,
		&t.ExpiresAt,
		&t.CreatedAt,
		&t.DeletedAt,
	)

	if err != nil {
		return nil, handleNotFound(err, "токен не найден")
	}

	return &t, nil
}

func (u *userSessionRepository) DeleteByToken(ctx context.Context, token string) error {
	_, err := u.db.Exec(ctx, `UPDATE user_sessions SET deleted_at = NOW() WHERE token = $1 AND deleted_at IS NULL`, token)
	return err
}

func (u *userSessionRepository) DeleteByUserId(ctx context.Context, userID int, tokenType domain.TokenType) error {
	_, err := u.db.Exec(ctx, `UPDATE user_sessions SET deleted_at = NOW() WHERE user_id = $1 AND type = $2 AND deleted_at IS NULL`, userID, tokenType)
	return err
}

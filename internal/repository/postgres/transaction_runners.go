package postgres

import (
	"context"

	"github.com/magomedcoder/gen/internal/domain"
	"gorm.io/gorm"
)

type chatTransactionRunner struct {
	db      *gorm.DB
	runners domain.RunnerRepository
}

func NewChatTransactionRunner(db *gorm.DB, runners domain.RunnerRepository) domain.ChatTransactionRunner {
	return &chatTransactionRunner{
		db:      db,
		runners: runners,
	}
}

func (t *chatTransactionRunner) WithinTx(ctx context.Context, fn func(ctx context.Context, r domain.ChatRepos) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, NewChatRepos(tx, t.runners))
	})
}

type authTransactionRunner struct{ db *gorm.DB }

func NewAuthTransactionRunner(db *gorm.DB) domain.AuthTransactionRunner {
	return &authTransactionRunner{db: db}
}

func (t *authTransactionRunner) WithinTx(ctx context.Context, fn func(ctx context.Context, r domain.AuthRepos) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, NewAuthRepos(tx))
	})
}

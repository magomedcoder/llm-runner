package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/magomedcoder/gen/internal/domain"
	"github.com/magomedcoder/gen/internal/repository/postgres/model"
	"gorm.io/gorm"
)

type chatPreferenceRepository struct {
	db      *gorm.DB
	runners domain.RunnerRepository
}

func NewChatPreferenceRepository(db *gorm.DB, runners domain.RunnerRepository) domain.ChatPreferenceRepository {
	return &chatPreferenceRepository{db: db, runners: runners}
}

func (r *chatPreferenceRepository) GetSelectedRunner(ctx context.Context, userID int) (string, error) {
	var rows []struct {
		Host string
		Port int32
	}
	err := r.db.WithContext(ctx).Table("chats").
		Select("runners.host", "runners.port").
		Joins("INNER JOIN runners ON runners.id = chats.selected_runner_id").
		Where("chats.user_id = ?", userID).
		Order("chats.updated_at DESC").
		Limit(1).
		Scan(&rows).Error
	if err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return "", nil
	}
	return domain.RunnerListenAddress(rows[0].Host, rows[0].Port), nil
}

func (r *chatPreferenceRepository) SetSelectedRunner(ctx context.Context, userID int, runner string) error {
	runner = strings.TrimSpace(runner)
	var id *int64
	if runner != "" {
		rid, ok, err := r.runners.FindIDByListenAddress(ctx, runner)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("раннер не найден для указанного адреса")
		}
		id = &rid
	}
	return r.db.WithContext(ctx).Model(&model.Chat{}).
		Where("user_id = ?", userID).
		Updates(map[string]any{
			"selected_runner_id": id,
			"updated_at":         gorm.Expr("NOW()"),
		}).Error
}

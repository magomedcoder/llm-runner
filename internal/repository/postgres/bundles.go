package postgres

import (
	"github.com/magomedcoder/gen/internal/domain"
	"gorm.io/gorm"
)

func NewChatRepos(db *gorm.DB, runners domain.RunnerRepository) domain.ChatRepos {
	return domain.ChatRepos{
		Session:         NewChatSessionRepository(db),
		Preference:      NewChatPreferenceRepository(db, runners),
		SessionSettings: NewChatSessionSettingsRepository(db),
		Message:         NewMessageRepository(db),
		MessageEdit:     NewMessageEditRepository(db),
		AssistantRegen:  NewAssistantMessageRegenerationRepository(db),
		File:            NewFileRepository(db),
	}
}

func NewAuthRepos(db *gorm.DB) domain.AuthRepos {
	return domain.AuthRepos{
		User:  NewUserRepository(db),
		Token: NewUserSessionRepository(db),
	}
}

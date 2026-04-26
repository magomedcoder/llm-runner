package model

import "time"

type EditorTextHistory struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int       `gorm:"column:user_id"`
	RunnerID  *int64    `gorm:"column:runner_id"`
	Text      string    `gorm:"column:text"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (EditorTextHistory) TableName() string {
	return "editor_text_history"
}

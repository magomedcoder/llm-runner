package model

import "time"

type MCPServer struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID         *int      `gorm:"column:user_id;index"`
	Name           string    `gorm:"column:name"`
	Enabled        bool      `gorm:"column:enabled"`
	Transport      string    `gorm:"column:transport"`
	Command        string    `gorm:"column:command"`
	ArgsJSON       string    `gorm:"column:args_json"`
	EnvJSON        string    `gorm:"column:env_json"`
	URL            string    `gorm:"column:url"`
	HeadersJSON    string    `gorm:"column:headers_json"`
	TimeoutSeconds int32     `gorm:"column:timeout_seconds"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (MCPServer) TableName() string {
	return "mcp_servers"
}

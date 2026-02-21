package config

import (
	"time"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
	LLMRunner   LLMRunnerConfig
	Runners     RunnersConfig
	Attachments AttachmentsConfig
}

type AttachmentsConfig struct {
	SaveDir string
}

type RunnersConfig struct {
	Addresses []string
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	DSN string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type LLMRunnerConfig struct {
	Address string
	Model   string
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: "50051",
			Host: "0.0.0.0",
		},
		Database: DatabaseConfig{
			DSN: "postgres://postgres:postgres@localhost:5432/gen_db?sslmode=disable",
		},
		JWT: JWTConfig{
			AccessSecret:  "gen",
			RefreshSecret: "gen",
			AccessTTL:     15 * time.Minute,
			RefreshTTL:    7 * 24 * time.Hour,
		},
		LLMRunner: LLMRunnerConfig{
			Address: "localhost:50052",
			Model:   "default",
		},
		Runners: RunnersConfig{
			Addresses: []string{},
		},
		Attachments: AttachmentsConfig{
			SaveDir: "./uploads",
		},
	}

	return config, nil
}

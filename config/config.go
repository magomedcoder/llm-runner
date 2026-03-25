package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type HostPort struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Config struct {
	Listen                   HostPort `yaml:"listen"`
	Core                     HostPort `yaml:"core"`
	RegistrationToken        string   `yaml:"registration_token"`
	LogLevel                 string   `yaml:"log_level"`
	ModelPath                string   `yaml:"model_path"`
	MmprojPath               string   `yaml:"mmproj_path"`
	DefaultModel             string   `yaml:"default_model"`
	MaxContextTokens         int      `yaml:"max_context_tokens"`
	MaxConcurrentGenerations int      `yaml:"max_concurrent_generations"`
	ModelRetention           string   `yaml:"model_retention"`
}

func Load() (*Config, error) {
	c := &Config{}

	configPath := os.Getenv("LLM_RUNNER_CONFIG")
	if configPath == "" {
		configPath = "./config.yaml"
	}

	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения конфигурационного файла %s: %w", configPath, err)
		}

		if err := yaml.Unmarshal(data, c); err != nil {
			return nil, fmt.Errorf("ошибка парсинга конфигурационного файла %s: %w", configPath, err)
		}
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) ListenAddr() string {
	h := strings.TrimSpace(c.Listen.Host)
	if h == "" || c.Listen.Port <= 0 {
		return ""
	}
	return net.JoinHostPort(h, strconv.Itoa(c.Listen.Port))
}

func (c *Config) CoreAddr() string {
	h := strings.TrimSpace(c.Core.Host)
	if h == "" || c.Core.Port <= 0 {
		return ""
	}
	return net.JoinHostPort(h, strconv.Itoa(c.Core.Port))
}

func (c *Config) validate() error {
	coreHost := strings.TrimSpace(c.Core.Host)
	hasCoreHost := coreHost != ""
	hasCorePort := c.Core.Port >= 1 && c.Core.Port <= 65535
	if hasCoreHost != hasCorePort {
		if hasCoreHost {
			return fmt.Errorf("core.port: укажите порт ядра от 1 до 65535")
		}
		return fmt.Errorf("core.host: укажите хост ядра")
	}

	if strings.TrimSpace(c.Listen.Host) == "" {
		return fmt.Errorf("listen.host: укажите хост для прослушивания")
	}
	if c.Listen.Port < 1 || c.Listen.Port > 65535 {
		return fmt.Errorf("listen.port: укажите порт от 1 до 65535")
	}
	if strings.TrimSpace(c.DefaultModel) == "" {
		return fmt.Errorf("default_model: укажите модель по умолчанию (имя из каталога model_path)")
	}
	r := strings.TrimSpace(c.ModelRetention)
	if r == "" {
		return nil
	}
	rl := strings.ToLower(r)
	if rl == "keep" || rl == "unload_after_rpc" {
		return nil
	}
	return fmt.Errorf("model_retention: неизвестное значение %q (допустимы keep, unload_after_rpc)", r)
}

func (c *Config) UnloadAfterRPC() bool {
	return strings.ToLower(strings.TrimSpace(c.ModelRetention)) == "unload_after_rpc"
}

func (c *Config) ModelsDir() string {
	p := strings.TrimSpace(c.ModelPath)
	if p == "" {
		return ""
	}

	if st, err := os.Stat(p); err == nil && !st.IsDir() {
		return filepath.Dir(p)
	}

	if strings.EqualFold(filepath.Ext(p), ".gguf") {
		return filepath.Dir(p)
	}

	return p
}

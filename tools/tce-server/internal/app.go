package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/magomedcoder/gen/internal/service"
)

type App struct {
	llm        *service.LLMRunnerService
	model      string
	runnerAddr string
	host       string
}

func NewFromEnv() (*App, error) {
	runnerAddr := strings.TrimSpace(os.Getenv("GEN_RUNNER_ADDR"))
	model := strings.TrimSpace(os.Getenv("GEN_MODEL"))
	host := strings.TrimSpace(os.Getenv("TCE_SERVER_HOST"))

	if runnerAddr == "" {
		runnerAddr = "0.0.0.0:50052"
	}

	if model == "" {
		model = "Qwen3-8B-Q8_0"
	}

	if host == "" {
		host = "127.0.0.1:8000"
	}

	llmSvc, err := service.NewLLMRunnerService(runnerAddr, model)
	if err != nil {
		return nil, fmt.Errorf("не удалось инициализировать клиент gen runner: %w", err)
	}

	return &App{
		llm:        llmSvc,
		model:      model,
		runnerAddr: runnerAddr,
		host:       host,
	}, nil
}

func (a *App) Close() {
	if a == nil || a.llm == nil {
		return
	}

	if err := a.llm.Close(); err != nil {
		log.Printf("предупреждение: не удалось закрыть клиент runner: %v", err)
	}
}

func (a *App) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/tce/v1/health", a.handleHealth)
	mux.HandleFunc("/tce/v1/chat", a.handleChat)
	mux.HandleFunc("/tce/v1/agent/step", a.handleAgentStep)

	log.Printf("Tce-server запущен на %s", a.host)
	log.Printf("раннер: %s, модель: %s", a.runnerAddr, a.model)
	log.Printf("проверка здоровья: GET  http://%s/tce/v1/health", a.host)
	log.Printf("чат:               POST http://%s/tce/v1/chat", a.host)
	log.Printf("шаг агента:        POST http://%s/tce/v1/agent/step", a.host)

	return http.ListenAndServe(a.host, withCORS(mux))
}

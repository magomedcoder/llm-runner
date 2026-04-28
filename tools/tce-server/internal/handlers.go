package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/magomedcoder/gen/internal/domain"
)

func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, "GET")
		return
	}

	ok, err := a.llm.CheckConnection(r.Context())
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, errorBody("internal_error", "раннер недоступен"))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok": ok,
	})
}

func (a *App) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, "POST")
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeBadRequest(w, "некорректное тело JSON")
		return
	}

	if req.Stream == nil {
		writeBadRequest(w, "поле stream обязательно")
		return
	}
	if req.System == nil {
		writeBadRequest(w, "поле system обязательно")
		return
	}
	if len(req.Messages) == 0 {
		writeBadRequest(w, "массив messages не должен быть пустым")
		return
	}
	if req.Messages[len(req.Messages)-1].Role != "user" {
		writeBadRequest(w, `последнее сообщение в messages должно иметь role "user"`)
		return
	}

	messages := makeRunnerMessages(*req.System, req.Messages, req.Editor)
	genParams := mapGenerateParams(req.Generate)
	ch, err := a.llm.SendMessage(r.Context(), 0, a.model, messages, nil, 0, genParams)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorBody("internal_error", "ошибка генерации"))
		return
	}

	if *req.Stream {
		writeRunnerSSE(w, ch)
		return
	}

	var full strings.Builder
	for chunk := range ch {
		if chunk.Content != "" {
			full.WriteString(chunk.Content)
		}
	}

	writeJSON(w, http.StatusOK, ChatResponse{
		Message: ChatMessage{
			Role:    "assistant",
			Content: full.String(),
		},
		Finish: "stop",
	})
}

func (a *App) handleAgentStep(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, "POST")
		return
	}

	if r.Body != nil && r.ContentLength != 0 {
		var req AgentStepRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeBadRequest(w, "некорректное тело JSON")
			return
		}
	}

	writeJSON(w, http.StatusOK, AgentStepResponse{
		Finish:  false,
		Summary: "Проверил файлы и подготовил следующие действия.",
		Calls: []AgentToolCall{
			{
				Tool: "read_file",
				ID:   "call-1",
				Args: map[string]any{
					"path": "src/main.rs",
				},
			},
		},
	})
}

func writeRunnerSSE(w http.ResponseWriter, chunks chan domain.LLMStreamChunk) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "сервер не поддерживает потоковую передачу",
		})
		return
	}

	for chunk := range chunks {
		if chunk.Content == "" {
			continue
		}

		fmt.Fprintf(w, "event: delta\n")
		fmt.Fprintf(w, "data: {\"text\":%q}\n\n", chunk.Content)
		flusher.Flush()
	}

	fmt.Fprintf(w, "event: end\n")
	fmt.Fprintf(w, "data: {\"finish\":\"stop\"}\n\n")
	flusher.Flush()
}

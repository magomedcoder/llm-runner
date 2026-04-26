package bitrix24server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/text/encoding/charmap"
)

func TestBitrixClientCall_DecodesWindows1251JSON(t *testing.T) {
	cp1251Body, err := charmap.Windows1251.NewEncoder().Bytes([]byte(`{"result":{"task":{"TITLE":"Привет"}}}`))
	if err != nil {
		t.Fatalf("encode cp1251: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=windows-1251")
		_, _ = w.Write(cp1251Body)
	}))
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, time.Second)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	resp, err := client.call(context.Background(), "tasks.task.get", map[string]any{"taskId": 1001})
	if err != nil {
		t.Fatalf("call: %v", err)
	}

	got := nestedTaskTitle(resp)
	if got != "Привет" {
		t.Fatalf("unexpected title: %q", got)
	}
}

func TestBitrixClientCall_DecodesWindows1251WithoutCharsetHint(t *testing.T) {
	cp1251Body, err := charmap.Windows1251.NewEncoder().Bytes([]byte(`{"result":{"task":{"TITLE":"Задача"}}}`))
	if err != nil {
		t.Fatalf("encode cp1251: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(cp1251Body)
	}))
	defer srv.Close()

	client, err := newBitrixClient(srv.URL, time.Second)
	if err != nil {
		t.Fatalf("newBitrixClient: %v", err)
	}

	resp, err := client.call(context.Background(), "tasks.task.get", map[string]any{"taskId": 1001})
	if err != nil {
		t.Fatalf("call: %v", err)
	}

	got := nestedTaskTitle(resp)
	if got != "Задача" {
		t.Fatalf("unexpected title: %q", got)
	}
}

func nestedTaskTitle(resp map[string]any) string {
	result, _ := resp["result"].(map[string]any)
	task, _ := result["task"].(map[string]any)
	title, _ := task["TITLE"].(string)
	return title
}

package bitrix24server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/magomedcoder/gen/pkg/mcpsafe"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Config struct {
	WebhookBase string
}

func NewServer(cfg Config) (*mcp.Server, error) {
	timeout := 20 * time.Second
	log.Printf("[b24-mcp] init server timeout=%s webhook_base_set=%t", timeout, strings.TrimSpace(cfg.WebhookBase) != "")

	client, err := newBitrixClient(cfg.WebhookBase, timeout)
	if err != nil {
		return nil, err
	}

	srv := mcp.NewServer(&mcp.Implementation{
		Name:    "bitrix24",
		Version: "0.1.0",
	}, nil)

	type callArgs struct {
		Method string         `json:"method" jsonschema:"REST-метод Bitrix24, например tasks.task.list"`
		Params map[string]any `json:"params" jsonschema:"Параметры метода в JSON"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "b24_call_method",
		Description: "Универсальный вызов метода Bitrix24 REST API",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args callArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-bitrix24", "b24_call_method", func() (*mcp.CallToolResult, any, error) {
			method := normalizeMethod(args.Method)
			if method == "" {
				return nil, nil, fmt.Errorf("method is required")
			}
			log.Printf("[b24-mcp] tool=b24_call_method method=%q params_keys=%d", method, len(args.Params))

			resp, err := client.call(ctx, method, args.Params)
			if err != nil {
				log.Printf("[b24-mcp] tool=b24_call_method method=%q err=%v", method, err)
				return nil, nil, err
			}
			log.Printf("[b24-mcp] tool=b24_call_method method=%q ok", method)

			return textResult(prettyJSON(resp)), nil, nil
		})
	})

	type listTasksArgs struct {
		Filter map[string]any `json:"filter,omitempty" jsonschema:"Фильтры задач Bitrix24 (например UF_* поля)"`
		Select []string       `json:"select,omitempty" jsonschema:"Список полей для выборки"`
		Order  map[string]any `json:"order,omitempty" jsonschema:"Сортировка, например {\"CREATED_DATE\":\"desc\"}"`
		Start  *int           `json:"start,omitempty" jsonschema:"Пагинация Bitrix24: offset"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "b24_list_tasks",
		Description: "Получить список задач Bitrix24 (tasks.task.list)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args listTasksArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-bitrix24", "b24_list_tasks", func() (*mcp.CallToolResult, any, error) {
			log.Printf("[b24-mcp] tool=b24_list_tasks filter_keys=%d select=%d order_keys=%d has_start=%t", len(args.Filter), len(args.Select), len(args.Order), args.Start != nil)
			payload := map[string]any{}
			if args.Filter != nil {
				payload["filter"] = args.Filter
			}

			if len(args.Select) > 0 {
				payload["select"] = args.Select
			}

			if args.Order != nil {
				payload["order"] = args.Order
			}

			if args.Start != nil {
				payload["start"] = *args.Start
			}

			resp, err := client.call(ctx, "tasks.task.list", payload)
			if err != nil {
				log.Printf("[b24-mcp] tool=b24_list_tasks err=%v", err)
				return nil, nil, err
			}
			log.Printf("[b24-mcp] tool=b24_list_tasks ok")

			return textResult(prettyJSON(resp)), nil, nil
		})
	})

	type getTaskArgs struct {
		TaskID int      `json:"task_id" jsonschema:"ID задачи"`
		Select []string `json:"select,omitempty" jsonschema:"Список полей для выборки"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "b24_get_task",
		Description: "Получить задачу по ID (tasks.task.get)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args getTaskArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-bitrix24", "b24_get_task", func() (*mcp.CallToolResult, any, error) {
			if args.TaskID <= 0 {
				return nil, nil, fmt.Errorf("task_id must be > 0")
			}
			log.Printf("[b24-mcp] tool=b24_get_task task_id=%d select=%d", args.TaskID, len(args.Select))

			payload := map[string]any{
				"taskId": args.TaskID,
			}

			if len(args.Select) > 0 {
				payload["select"] = args.Select
			}

			resp, err := client.call(ctx, "tasks.task.get", payload)
			if err != nil {
				log.Printf("[b24-mcp] tool=b24_get_task task_id=%d err=%v", args.TaskID, err)
				return nil, nil, err
			}
			log.Printf("[b24-mcp] tool=b24_get_task task_id=%d ok", args.TaskID)

			return textResult(prettyJSON(resp)), nil, nil
		})
	})

	type getCommentsArgs struct {
		TaskID int            `json:"task_id" jsonschema:"ID задачи"`
		Order  map[string]any `json:"order,omitempty" jsonschema:"Сортировка комментариев, например {\"ID\":\"asc\"}"`
		Select []string       `json:"select,omitempty" jsonschema:"Поля комментариев"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "b24_get_task_comments",
		Description: "Получить комментарии задачи (tasks.task.commentitem.getlist)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args getCommentsArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-bitrix24", "b24_get_task_comments", func() (*mcp.CallToolResult, any, error) {
			if args.TaskID <= 0 {
				return nil, nil, fmt.Errorf("task_id must be > 0")
			}
			log.Printf("[b24-mcp] tool=b24_get_task_comments task_id=%d select=%d order_keys=%d", args.TaskID, len(args.Select), len(args.Order))

			payload := map[string]any{
				"taskId": args.TaskID,
			}

			if args.Order != nil {
				payload["order"] = args.Order
			}

			if len(args.Select) > 0 {
				payload["select"] = args.Select
			}

			resp, err := client.call(ctx, "tasks.task.commentitem.getlist", payload)
			if err != nil {
				log.Printf("[b24-mcp] tool=b24_get_task_comments task_id=%d err=%v", args.TaskID, err)
				return nil, nil, err
			}
			log.Printf("[b24-mcp] tool=b24_get_task_comments task_id=%d ok", args.TaskID)

			return textResult(prettyJSON(resp)), nil, nil
		})
	})

	type analyzeArgs struct {
		TaskID          int   `json:"task_id" jsonschema:"ID задачи для анализа"`
		IncludeComments *bool `json:"include_comments,omitempty" jsonschema:"Подтянуть комментарии для анализа (по умолчанию true)"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "b24_analyze_task",
		Description: "Быстрый анализ задачи: статус, дедлайн, активность, риски",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args analyzeArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-bitrix24", "b24_analyze_task", func() (*mcp.CallToolResult, any, error) {
			if args.TaskID <= 0 {
				return nil, nil, fmt.Errorf("task_id must be > 0")
			}

			includeComments := true
			if args.IncludeComments != nil {
				includeComments = *args.IncludeComments
			}
			log.Printf("[b24-mcp] tool=b24_analyze_task task_id=%d include_comments=%t", args.TaskID, includeComments)

			taskResp, err := client.call(ctx, "tasks.task.get", map[string]any{
				"taskId": args.TaskID,
				"select": []string{
					"ID",
					"TITLE",
					"CREATED_DATE",
					"DEADLINE",
					"STATUS",
					"CREATED_BY",
					"RESPONSIBLE_ID",
				},
			})
			if err != nil {
				log.Printf("[b24-mcp] tool=b24_analyze_task task_id=%d stage=task_get err=%v", args.TaskID, err)
				return nil, nil, err
			}

			task, err := extractTask(taskResp)
			if err != nil {
				log.Printf("[b24-mcp] tool=b24_analyze_task task_id=%d stage=extract_task err=%v", args.TaskID, err)
				return nil, nil, err
			}

			comments := []map[string]any{}
			if includeComments {
				commentResp, err := client.call(ctx, "tasks.task.commentitem.getlist", map[string]any{
					"taskId": args.TaskID,
					"order":  map[string]string{"ID": "asc"},
				})

				if err != nil {
					log.Printf("[b24-mcp] tool=b24_analyze_task task_id=%d stage=comments err=%v", args.TaskID, err)
					return nil, nil, err
				}

				comments = extractComments(commentResp)
			}

			analysis := analyzeTask(task, comments, time.Now())
			log.Printf("[b24-mcp] tool=b24_analyze_task task_id=%d ok comments=%d", args.TaskID, len(comments))
			return textResult(analysis), nil, nil
		})
	})

	return srv, nil
}

func normalizeMethod(method string) string {
	method = strings.TrimSpace(method)
	method = strings.TrimPrefix(method, "/")
	method = strings.TrimSuffix(method, "/")
	method = strings.TrimSuffix(method, ".json")
	return method
}

func prettyJSON(data any) string {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", data)
	}

	return string(b)
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}
}

func extractTask(resp map[string]any) (map[string]any, error) {
	result, ok := resp["result"]
	if !ok {
		return nil, fmt.Errorf("bitrix response has no result field")
	}

	asMap, ok := result.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("bitrix result is not an object")
	}

	taskRaw, ok := asMap["task"]
	if !ok {
		return nil, fmt.Errorf("bitrix result.task is missing")
	}

	task, ok := taskRaw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("bitrix result.task has unexpected type")
	}
	return task, nil
}

func extractComments(resp map[string]any) []map[string]any {
	result, ok := resp["result"]
	if !ok {
		return nil
	}

	switch typed := result.(type) {
	case []any:
		return toMapSlice(typed)
	case map[string]any:
		if comments, ok := typed["comments"]; ok {
			if arr, ok := comments.([]any); ok {
				return toMapSlice(arr)
			}
		}
		if arr, ok := typed["items"].([]any); ok {
			return toMapSlice(arr)
		}
	}

	return nil
}

func toMapSlice(items []any) []map[string]any {
	out := make([]map[string]any, 0, len(items))
	for _, item := range items {
		asMap, ok := item.(map[string]any)
		if ok {
			out = append(out, asMap)
		}
	}

	return out
}

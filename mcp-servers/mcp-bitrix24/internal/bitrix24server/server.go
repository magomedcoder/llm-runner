package bitrix24server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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
		Version: "1.0.0",
	}, nil)

	type listTasksArgs struct {
		Filter map[string]any `json:"filter,omitempty" jsonschema:"Фильтры задач Bitrix24 (например UF_* поля)"`
		Select []string       `json:"select,omitempty" jsonschema:"Список полей для выборки"`
		Order  map[string]any `json:"order,omitempty" jsonschema:"Сортировка, например {\"CREATED_DATE\":\"desc\"}"`
		Params map[string]any `json:"params,omitempty" jsonschema:"Дополнительные параметры, например {\"WITH_TIMER_INFO\":true}"`
		Start  *int           `json:"start,omitempty" jsonschema:"Пагинация Bitrix24: offset"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "b24_list_tasks",
		Description: "Получить список задач Bitrix24 (tasks.task.list)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args listTasksArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-bitrix24", "b24_list_tasks", func() (*mcp.CallToolResult, any, error) {
			logToolArgs("b24_list_tasks", args)
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

			if args.Params != nil {
				payload["params"] = args.Params
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
			logToolResult("b24_list_tasks", resp)

			return textResult(prettyJSON(resp)), nil, nil
		})
	})

	type getTaskArgs struct {
		TaskID any      `json:"task_id" jsonschema:"ID задачи (число или строка с цифрами)"`
		Select []string `json:"select,omitempty" jsonschema:"Список полей для выборки"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "b24_get_task",
		Description: "Получить задачу по ID (tasks.task.get)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args getTaskArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-bitrix24", "b24_get_task", func() (*mcp.CallToolResult, any, error) {
			logToolArgs("b24_get_task", args)
			taskID, err := parseTaskID(args.TaskID)
			if err != nil {
				return nil, nil, err
			}
			if taskID <= 0 {
				return nil, nil, fmt.Errorf("task_id must be > 0")
			}
			log.Printf("[b24-mcp] tool=b24_get_task task_id=%d select=%d", taskID, len(args.Select))

			payload := map[string]any{
				"taskId": taskID,
			}

			if len(args.Select) > 0 {
				payload["select"] = args.Select
			}

			resp, err := client.call(ctx, "tasks.task.get", payload)
			if err != nil {
				log.Printf("[b24-mcp] tool=b24_get_task task_id=%d err=%v", taskID, err)
				return nil, nil, err
			}
			log.Printf("[b24-mcp] tool=b24_get_task task_id=%d ok", taskID)
			logToolResult("b24_get_task", resp)

			return textResult(prettyJSON(resp)), nil, nil
		})
	})

	type getCommentsArgs struct {
		TaskID any            `json:"task_id" jsonschema:"ID задачи (число или строка с цифрами)"`
		Order  map[string]any `json:"order,omitempty" jsonschema:"Сортировка комментариев, например {\"POST_DATE\":\"asc\"}"`
		Filter map[string]any `json:"filter,omitempty" jsonschema:"Фильтр комментариев, например {\"AUTHOR_ID\":503}"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "b24_get_task_comments",
		Description: "Получить комментарии задачи (task.commentitem.getlist)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args getCommentsArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-bitrix24", "b24_get_task_comments", func() (*mcp.CallToolResult, any, error) {
			logToolArgs("b24_get_task_comments", args)
			taskID, err := parseTaskID(args.TaskID)
			if err != nil {
				return nil, nil, err
			}
			if taskID <= 0 {
				return nil, nil, fmt.Errorf("task_id must be > 0")
			}
			log.Printf("[b24-mcp] tool=b24_get_task_comments task_id=%d order_keys=%d filter_keys=%d", taskID, len(args.Order), len(args.Filter))

			resp, err := client.callTaskCommentItemGetList(ctx, taskID, args.Order, args.Filter)
			if err != nil {
				log.Printf("[b24-mcp] tool=b24_get_task_comments task_id=%d err=%v", taskID, err)
				return nil, nil, err
			}
			log.Printf("[b24-mcp] tool=b24_get_task_comments task_id=%d ok", taskID)
			logToolResult("b24_get_task_comments", resp)

			return textResult(prettyJSON(resp)), nil, nil
		})
	})

	type analyzeArgs struct {
		TaskID          any   `json:"task_id,omitempty" jsonschema:"ID задачи для анализа (число или строка с цифрами)"`
		IncludeComments *bool `json:"include_comments,omitempty" jsonschema:"Подтянуть комментарии для анализа (по умолчанию true)"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "b24_analyze_task",
		Description: "Быстрый анализ задачи: статус, дедлайн, активность, риски",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args analyzeArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-bitrix24", "b24_analyze_task", func() (*mcp.CallToolResult, any, error) {
			logToolArgs("b24_analyze_task", args)
			taskIDPtr, err := parseOptionalTaskID(args.TaskID)
			if err != nil {
				return nil, nil, err
			}
			if taskIDPtr == nil {
				msg := "Для b24_analyze_task нужно передать task_id. " +
					"Пример: {\"task_id\": 1822404}. " +
					"Если нужен обзор без конкретной задачи, используйте b24_analyze_tasks_by_query."
				log.Printf("[b24-mcp] tool=b24_analyze_task missing_task_id")
				logToolResult("b24_analyze_task", msg)
				return textResult(msg), nil, nil
			}
			taskID := *taskIDPtr
			if taskID <= 0 {
				return nil, nil, fmt.Errorf("task_id must be > 0")
			}

			includeComments := true
			if args.IncludeComments != nil {
				includeComments = *args.IncludeComments
			}

			log.Printf("[b24-mcp] tool=b24_analyze_task task_id=%d include_comments=%t", taskID, includeComments)

			taskResp, err := client.call(ctx, "tasks.task.get", map[string]any{
				"taskId": taskID,
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
				log.Printf("[b24-mcp] tool=b24_analyze_task task_id=%d stage=task_get err=%v", taskID, err)
				return nil, nil, err
			}

			task, err := extractTask(taskResp)
			if err != nil {
				log.Printf("[b24-mcp] tool=b24_analyze_task task_id=%d stage=extract_task err=%v", taskID, err)
				return nil, nil, err
			}

			var comments []map[string]any
			if includeComments {
				commentResp, err := client.callTaskCommentItemGetList(ctx, taskID, map[string]any{"ID": "asc"}, nil)

				if err != nil {
					if isIgnorableCommentError(err) {
						log.Printf("[b24-mcp] tool=b24_analyze_task task_id=%d stage=comments soft_skip err=%v", taskID, err)
					} else {
						log.Printf("[b24-mcp] tool=b24_analyze_task task_id=%d stage=comments err=%v", taskID, err)
						return nil, nil, err
					}
				} else {
					comments = extractComments(commentResp)
				}
			}

			analysis := analyzeTask(task, comments, time.Now())
			log.Printf("[b24-mcp] tool=b24_analyze_task task_id=%d ok comments=%d", taskID, len(comments))
			logToolResult("b24_analyze_task", analysis)
			return textResult(analysis), nil, nil
		})
	})

	type analyticsQueryArgs struct {
		Query           string         `json:"query,omitempty" jsonschema:"Текст запроса по аналитике, например: 'покажи просроченные и блокеры'. Если не задан, будет общий аналитический обзор."`
		TaskID          any            `json:"task_id,omitempty" jsonschema:"Опционально: анализ конкретной задачи (число или строка с цифрами)"`
		Filter          map[string]any `json:"filter,omitempty" jsonschema:"Фильтр для tasks.task.list, если нужен анализ набора задач"`
		Order           map[string]any `json:"order,omitempty" jsonschema:"Сортировка для tasks.task.list"`
		Start           *int           `json:"start,omitempty" jsonschema:"Смещение в tasks.task.list"`
		Limit           *int           `json:"limit,omitempty" jsonschema:"Макс. задач для аналитики (1..50), по умолчанию 20"`
		IncludeComments *bool          `json:"include_comments,omitempty" jsonschema:"Подтягивать комментарии при аналитике (по умолчанию false)"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "b24_analyze_tasks_by_query",
		Description: "Выполнить пользовательский запрос по аналитике задач (просрочки, риски, блокеры, активность)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args analyticsQueryArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-bitrix24", "b24_analyze_tasks_by_query", func() (*mcp.CallToolResult, any, error) {
			logToolArgs("b24_analyze_tasks_by_query", args)
			query := strings.TrimSpace(args.Query)
			if query == "" {
				query = "общий аналитический обзор"
			}
			taskID, err := parseOptionalTaskID(args.TaskID)
			if err != nil {
				return nil, nil, err
			}
			result, err := runAnalyticsQuery(ctx, client, query, taskID, args.Filter, args.Order, args.Start, args.Limit, args.IncludeComments)
			if err != nil {
				return nil, nil, err
			}
			logToolResult("b24_analyze_tasks_by_query", result)
			return textResult(result), nil, nil
		})
	})

	type portfolioArgs struct {
		Filter          map[string]any `json:"filter,omitempty" jsonschema:"Фильтр задач для выборки портфеля"`
		Order           map[string]any `json:"order,omitempty" jsonschema:"Сортировка tasks.task.list"`
		Start           *int           `json:"start,omitempty" jsonschema:"Смещение в tasks.task.list"`
		Limit           *int           `json:"limit,omitempty" jsonschema:"Макс. задач для анализа (1..50), по умолчанию 30"`
		IncludeComments *bool          `json:"include_comments,omitempty" jsonschema:"Подтягивать комментарии (по умолчанию false)"`
		GroupBy         string         `json:"group_by,omitempty" jsonschema:"Группировка: responsible|creator|status (по умолчанию responsible)"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "b24_analyze_tasks_portfolio",
		Description: "Портфельная аналитика задач: сводка по ответственным/постановщикам, просрочки и риски",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args portfolioArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-bitrix24", "b24_analyze_tasks_portfolio", func() (*mcp.CallToolResult, any, error) {
			logToolArgs("b24_analyze_tasks_portfolio", args)
			result, err := runPortfolioAnalytics(ctx, client, args.Filter, args.Order, args.Start, args.Limit, args.IncludeComments, args.GroupBy)
			if err != nil {
				return nil, nil, err
			}
			logToolResult("b24_analyze_tasks_portfolio", result)
			return textResult(result), nil, nil
		})
	})

	type executiveSummaryArgs struct {
		Filter          map[string]any `json:"filter,omitempty" jsonschema:"Фильтр задач для выборки отчета"`
		Order           map[string]any `json:"order,omitempty" jsonschema:"Сортировка tasks.task.list"`
		Start           *int           `json:"start,omitempty" jsonschema:"Смещение в tasks.task.list"`
		Limit           *int           `json:"limit,omitempty" jsonschema:"Макс. задач для отчета (1..50), по умолчанию 40"`
		PeriodDays      *int           `json:"period_days,omitempty" jsonschema:"Длина периода в днях для тренда (1..30), по умолчанию 7"`
		IncludeComments *bool          `json:"include_comments,omitempty" jsonschema:"Подтягивать комментарии для оценки блокеров (по умолчанию false)"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "b24_analyze_tasks_executive_summary",
		Description: "Управленческая сводка по задачам за период с трендами и ключевыми рисками",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args executiveSummaryArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-bitrix24", "b24_analyze_tasks_executive_summary", func() (*mcp.CallToolResult, any, error) {
			logToolArgs("b24_analyze_tasks_executive_summary", args)
			result, err := runExecutiveSummary(ctx, client, args.Filter, args.Order, args.Start, args.Limit, args.PeriodDays, args.IncludeComments)
			if err != nil {
				return nil, nil, err
			}
			logToolResult("b24_analyze_tasks_executive_summary", result)
			return textResult(result), nil, nil
		})
	})

	type slaArgs struct {
		Filter             map[string]any `json:"filter,omitempty" jsonschema:"Фильтр задач для SLA-аналитики"`
		Order              map[string]any `json:"order,omitempty" jsonschema:"Сортировка tasks.task.list"`
		Start              *int           `json:"start,omitempty" jsonschema:"Смещение в tasks.task.list"`
		Limit              *int           `json:"limit,omitempty" jsonschema:"Макс. задач (1..50), по умолчанию 40"`
		SoonHoursThreshold *int           `json:"soon_hours_threshold,omitempty" jsonschema:"Порог 'скоро дедлайн' в часах, по умолчанию 24"`
		IncludeComments    *bool          `json:"include_comments,omitempty" jsonschema:"Подтягивать комментарии для блокеров (по умолчанию false)"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "b24_analyze_tasks_sla",
		Description: "SLA-контроль задач: просрочки, дедлайн-сегодня/скоро, приоритет реакции",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, args slaArgs) (*mcp.CallToolResult, any, error) {
		return mcpsafe.SafeToolInvoke("mcp-bitrix24", "b24_analyze_tasks_sla", func() (*mcp.CallToolResult, any, error) {
			logToolArgs("b24_analyze_tasks_sla", args)
			result, err := runSLASummary(ctx, client, args.Filter, args.Order, args.Start, args.Limit, args.SoonHoursThreshold, args.IncludeComments)
			if err != nil {
				return nil, nil, err
			}
			logToolResult("b24_analyze_tasks_sla", result)
			return textResult(result), nil, nil
		})
	})

	return srv, nil
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

func parseTaskID(v any) (int, error) {
	switch typed := v.(type) {
	case int:
		return typed, nil
	case int64:
		return int(typed), nil
	case float64:
		return int(typed), nil
	case string:
		s := strings.TrimSpace(typed)
		if s == "" {
			return 0, fmt.Errorf("task_id is required")
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("task_id must be integer-like, got %q", typed)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("task_id has unsupported type %T", v)
	}
}

func parseOptionalTaskID(v any) (*int, error) {
	if v == nil {
		return nil, nil
	}
	taskID, err := parseTaskID(v)
	if err != nil {
		return nil, err
	}
	return &taskID, nil
}

func logToolArgs(tool string, args any) {
	raw, err := json.Marshal(args)
	if err != nil {
		log.Printf("[b24-mcp] tool=%s args_marshal_err=%v", tool, err)
		return
	}
	log.Printf("[b24-mcp] tool=%s args=%s", tool, truncateForLog(prettyJSONString(raw), 3000))
}

func logToolResult(tool string, result any) {
	raw, err := json.Marshal(result)
	if err != nil {
		log.Printf("[b24-mcp] tool=%s result_marshal_err=%v", tool, err)
		return
	}
	log.Printf("[b24-mcp] tool=%s result=%s", tool, truncateForLog(prettyJSONString(raw), 3000))
}

package template

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"math"
	"slices"
	"strings"
	"sync"
	texttpl "text/template"
	"text/template/parse"
	"time"
)

var ErrNoMatchingPreset = errors.New("шаблон: нет пресета, достаточно близкого к chat_template модели")

var templatesOnce = sync.OnceValues(loadAllPresets)

func loadAllPresets() ([]*MatchedPreset, error) {
	templates, err := loadPresetEntries()
	if err != nil {
		return nil, err
	}
	if err := hydratePresetAssets(templates); err != nil {
		return nil, err
	}
	return templates, nil
}

type MatchedPreset struct {
	Name       string `json:"name"`
	IndexJinja string `json:"template"`
	Bytes      []byte `json:"-"`

	Parameters *struct {
		Stop []string `json:"stop"`
	}
}

func (t *MatchedPreset) Reader() io.Reader {
	return bytes.NewReader(t.Bytes)
}

func PresetStopSequences(p *MatchedPreset) []string {
	if p == nil || p.Parameters == nil {
		return nil
	}

	return p.Parameters.Stop
}

func presetMatchDistance(modelChatTemplate, presetFingerprint string) int {
	d := stringEditDistance(modelChatTemplate, presetFingerprint)
	m := strings.ToLower(modelChatTemplate)
	p := strings.ToLower(presetFingerprint)
	if strings.Contains(m, "<|im_start|>") && !strings.Contains(p, "im_start") {
		d += 900
	}
	if strings.Contains(m, "<|start_header_id|>") && !strings.Contains(p, "<|start_header_id|>") {
		d += 900
	}
	if strings.Contains(m, "[inst]") && !strings.Contains(p, "[inst]") {
		d += 900
	}
	return d
}

func Named(modelChatTemplateJinja string) (*MatchedPreset, error) {
	templates, err := templatesOnce()
	if err != nil {
		return nil, err
	}
	if len(templates) == 0 {
		return nil, errors.New("пресеты: список пуст (нет загруженных presets/*.json или все пресеты отфильтрованы)")
	}

	var best *MatchedPreset
	score := math.MaxInt
	for _, t := range templates {
		if d := presetMatchDistance(modelChatTemplateJinja, t.IndexJinja); d < score {
			score, best = d, t
		}
	}

	if score < maxPresetMatchDistance {
		return best, nil
	}

	hint := ""
	if best != nil {
		hint = fmt.Sprintf("расстояние=%d ближайший_пресет=%q", score, best.Name)
	}
	return nil, fmt.Errorf("%w (%s)", ErrNoMatchingPreset, hint)
}

type Template struct {
	*texttpl.Template
	raw string
}

var response = parse.ActionNode{
	NodeType: parse.NodeAction,
	Pipe: &parse.PipeNode{
		NodeType: parse.NodePipe,
		Cmds: []*parse.CommandNode{
			{
				NodeType: parse.NodeCommand,
				Args: []parse.Node{
					&parse.FieldNode{
						NodeType: parse.NodeField,
						Ident:    []string{"Response"},
					},
				},
			},
		},
	},
}

var funcs = texttpl.FuncMap{
	"json": func(v any) string {
		b, _ := json.Marshal(v)
		return string(b)
	},
	"currentDate": func(args ...string) string {
		return time.Now().Format("2006-01-02")
	},
	"yesterdayDate": func(args ...string) string {
		return time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	},
	"toTypeScriptType": func(v any) string {
		if param, ok := v.(ToolProperty); ok {
			return param.ToTypeScriptType()
		}
		if param, ok := v.(*ToolProperty); ok && param != nil {
			return param.ToTypeScriptType()
		}
		return "any"
	},
}

func Parse(s string) (*Template, error) {
	tmpl := texttpl.New("").Option("missingkey=zero").Funcs(funcs)

	tmpl, err := tmpl.Parse(s)
	if err != nil {
		return nil, err
	}

	t := Template{Template: tmpl, raw: s}
	vars, err := t.Vars()
	if err != nil {
		return nil, err
	}

	if !slices.Contains(vars, "messages") && !slices.Contains(vars, "response") {
		tmpl.Tree.Root.Nodes = append(tmpl.Tree.Root.Nodes, &response)
	}

	return &t, nil
}

func (t *Template) Vars() ([]string, error) {
	var vars []string
	for _, tt := range t.Templates() {
		for _, n := range tt.Root.Nodes {
			v, err := Identifiers(n)
			if err != nil {
				return vars, err
			}
			vars = append(vars, v...)
		}
	}

	set := make(map[string]struct{})
	for _, n := range vars {
		set[strings.ToLower(n)] = struct{}{}
	}

	return slices.Sorted(maps.Keys(set)), nil
}

type Values struct {
	Messages []Message
	Tools
	Prompt string
	Suffix string
	Think  bool

	ThinkLevel string
	IsThinkSet bool

	forceLegacy bool
}

func (t *Template) Execute(w io.Writer, v Values) error {
	system, messages := collate(v.Messages)
	vars, err := t.Vars()
	if err != nil {
		return err
	}
	if v.Prompt != "" && v.Suffix != "" {
		return t.Template.Execute(w, map[string]any{
			"Prompt":     v.Prompt,
			"Suffix":     v.Suffix,
			"Response":   "",
			"Think":      v.Think,
			"ThinkLevel": v.ThinkLevel,
			"IsThinkSet": v.IsThinkSet,
		})
	} else if !v.forceLegacy && slices.Contains(vars, "messages") {
		return t.Template.Execute(w, map[string]any{
			"System":     system,
			"Messages":   convertMessagesForTemplate(messages),
			"Tools":      convertToolsForTemplate(v.Tools),
			"Response":   "",
			"Think":      v.Think,
			"ThinkLevel": v.ThinkLevel,
			"IsThinkSet": v.IsThinkSet,
		})
	}

	system = ""
	var b bytes.Buffer
	var prompt, response string
	for _, m := range messages {
		execute := func() error {
			if err := t.Template.Execute(&b, map[string]any{
				"System":     system,
				"Prompt":     prompt,
				"Response":   response,
				"Think":      v.Think,
				"ThinkLevel": v.ThinkLevel,
				"IsThinkSet": v.IsThinkSet,
			}); err != nil {
				return err
			}

			system = ""
			prompt = ""
			response = ""
			return nil
		}

		switch m.Role {
		case "system":
			if prompt != "" || response != "" {
				if err := execute(); err != nil {
					return err
				}
			}
			system = m.Content
		case "user":
			if response != "" {
				if err := execute(); err != nil {
					return err
				}
			}
			prompt = m.Content
		case "assistant":
			response = m.Content
		}
	}

	var cut bool
	nodes := deleteNode(t.Template.Root.Copy(), func(n parse.Node) bool {
		if field, ok := n.(*parse.FieldNode); ok && slices.Contains(field.Ident, "Response") {
			cut = true
			return false
		}

		return cut
	})

	tree := parse.Tree{Root: nodes.(*parse.ListNode)}
	if err := texttpl.Must(texttpl.New("").AddParseTree("", &tree)).Execute(&b, map[string]any{
		"System":     system,
		"Prompt":     prompt,
		"Response":   response,
		"Think":      v.Think,
		"ThinkLevel": v.ThinkLevel,
		"IsThinkSet": v.IsThinkSet,
	}); err != nil {
		return err
	}

	_, err = io.Copy(w, &b)
	return err
}

func collate(msgs []Message) (string, []*Message) {
	var system []string
	var collated []*Message
	for i := range msgs {
		if msgs[i].Role == "system" {
			system = append(system, msgs[i].Content)
		}

		if len(collated) > 0 && collated[len(collated)-1].Role == msgs[i].Role && msgs[i].Role != "tool" {
			collated[len(collated)-1].Content += "\n\n" + msgs[i].Content
		} else {
			collated = append(collated, &msgs[i])
		}
	}

	return strings.Join(system, "\n\n"), collated
}

type templateTools []templateTool

func (t templateTools) String() string {
	bts, _ := json.Marshal(t)
	return string(bts)
}

type templateArgs map[string]any

func (t templateArgs) String() string {
	if t == nil {
		return "{}"
	}
	bts, _ := json.Marshal(t)
	return string(bts)
}

type templateProperties map[string]ToolProperty

func (t templateProperties) String() string {
	if t == nil {
		return "{}"
	}
	bts, _ := json.Marshal(t)
	return string(bts)
}

type templateTool struct {
	Type     string               `json:"type"`
	Items    any                  `json:"items,omitempty"`
	Function templateToolFunction `json:"function"`
}

type templateToolFunction struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	Parameters  templateToolFunctionParameters `json:"parameters"`
}

type templateToolFunctionParameters struct {
	Type       string             `json:"type"`
	Defs       any                `json:"$defs,omitempty"`
	Items      any                `json:"items,omitempty"`
	Required   []string           `json:"required,omitempty"`
	Properties templateProperties `json:"properties"`
}

type templateToolCall struct {
	ID       string
	Function templateToolCallFunction
}

type templateToolCallFunction struct {
	Index     int
	Name      string
	Arguments templateArgs
}

type templateMessage struct {
	Role       string
	Content    string
	Thinking   string
	Images     []ImageData
	ToolCalls  []templateToolCall
	ToolName   string
	ToolCallID string
}

func convertToolsForTemplate(tools Tools) templateTools {
	if tools == nil {
		return nil
	}
	result := make(templateTools, len(tools))
	for i, tool := range tools {
		var props templateProperties
		if tool.Function.Parameters.Properties != nil {
			props = templateProperties(tool.Function.Parameters.Properties.ToMap())
		}
		result[i] = templateTool{
			Type:  tool.Type,
			Items: tool.Items,
			Function: templateToolFunction{
				Name:        tool.Function.Name,
				Description: tool.Function.Description,
				Parameters: templateToolFunctionParameters{
					Type:       tool.Function.Parameters.Type,
					Defs:       tool.Function.Parameters.Defs,
					Items:      tool.Function.Parameters.Items,
					Required:   tool.Function.Parameters.Required,
					Properties: props,
				},
			},
		}
	}
	return result
}

func convertMessagesForTemplate(messages []*Message) []*templateMessage {
	if messages == nil {
		return nil
	}
	result := make([]*templateMessage, len(messages))
	for i, msg := range messages {
		var toolCalls []templateToolCall
		for _, tc := range msg.ToolCalls {
			toolCalls = append(toolCalls, templateToolCall{
				ID: tc.ID,
				Function: templateToolCallFunction{
					Index:     tc.Function.Index,
					Name:      tc.Function.Name,
					Arguments: templateArgs(tc.Function.Arguments.ToMap()),
				},
			})
		}
		result[i] = &templateMessage{
			Role:       msg.Role,
			Content:    msg.Content,
			Thinking:   msg.Thinking,
			Images:     msg.Images,
			ToolCalls:  toolCalls,
			ToolName:   msg.ToolName,
			ToolCallID: msg.ToolCallID,
		}
	}
	return result
}

func Identifiers(n parse.Node) ([]string, error) {
	switch n := n.(type) {
	case *parse.ListNode:
		var names []string
		for _, n := range n.Nodes {
			i, err := Identifiers(n)
			if err != nil {
				return names, err
			}
			names = append(names, i...)
		}

		return names, nil
	case *parse.TemplateNode:
		if n.Pipe == nil {
			return nil, errors.New("шаблон Go: в узле Template не задан pipe (некорректная разметка шаблона)")
		}
		return Identifiers(n.Pipe)
	case *parse.ActionNode:
		if n.Pipe == nil {
			return nil, errors.New("шаблон Go: в узле Action не задан pipe (некорректная разметка шаблона)")
		}
		return Identifiers(n.Pipe)
	case *parse.BranchNode:
		if n.Pipe == nil {
			return nil, errors.New("шаблон Go: в узле ветвления (if/range) не задан pipe (некорректная разметка шаблона)")
		}
		names, err := Identifiers(n.Pipe)
		if err != nil {
			return names, err
		}
		for _, n := range []*parse.ListNode{n.List, n.ElseList} {
			if n != nil {
				i, err := Identifiers(n)
				if err != nil {
					return names, err
				}
				names = append(names, i...)
			}
		}
		return names, nil
	case *parse.IfNode:
		return Identifiers(&n.BranchNode)
	case *parse.RangeNode:
		return Identifiers(&n.BranchNode)
	case *parse.WithNode:
		return Identifiers(&n.BranchNode)
	case *parse.PipeNode:
		var names []string
		for _, c := range n.Cmds {
			for _, a := range c.Args {
				i, err := Identifiers(a)
				if err != nil {
					return names, err
				}
				names = append(names, i...)
			}
		}
		return names, nil
	case *parse.FieldNode:
		return n.Ident, nil
	case *parse.VariableNode:
		return n.Ident, nil
	}

	return nil, nil
}

func deleteNode(n parse.Node, fn func(parse.Node) bool) parse.Node {
	var walk func(n parse.Node) parse.Node
	walk = func(n parse.Node) parse.Node {
		if fn(n) {
			return nil
		}

		switch t := n.(type) {
		case *parse.ListNode:
			var nodes []parse.Node
			for _, c := range t.Nodes {
				if n := walk(c); n != nil {
					nodes = append(nodes, n)
				}
			}

			t.Nodes = nodes
			return t
		case *parse.IfNode:
			t.BranchNode = *(walk(&t.BranchNode).(*parse.BranchNode))
		case *parse.WithNode:
			t.BranchNode = *(walk(&t.BranchNode).(*parse.BranchNode))
		case *parse.RangeNode:
			t.BranchNode = *(walk(&t.BranchNode).(*parse.BranchNode))
		case *parse.BranchNode:
			t.List = walk(t.List).(*parse.ListNode)
			if t.ElseList != nil {
				t.ElseList = walk(t.ElseList).(*parse.ListNode)
			}
		case *parse.ActionNode:
			n := walk(t.Pipe)
			if n == nil {
				return nil
			}

			t.Pipe = n.(*parse.PipeNode)
		case *parse.PipeNode:
			var commands []*parse.CommandNode
			for _, c := range t.Cmds {
				var args []parse.Node
				for _, a := range c.Args {
					if n := walk(a); n != nil {
						args = append(args, n)
					}
				}

				if len(args) == 0 {
					return nil
				}

				c.Args = args
				commands = append(commands, c)
			}

			if len(commands) == 0 {
				return nil
			}

			t.Cmds = commands
		}

		return n
	}

	return walk(n)
}

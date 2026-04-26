package template

import (
	"encoding/json"
	"fmt"
	"maps"
	"strings"
)

type ImageData []byte

type Message struct {
	Role       string      `json:"role"`
	Content    string      `json:"content"`
	Thinking   string      `json:"thinking,omitempty"`
	Images     []ImageData `json:"images,omitempty"`
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
	ToolName   string      `json:"tool_name,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
}

type Tools []Tool

type ToolCall struct {
	ID       string           `json:"id,omitempty"`
	Function ToolCallFunction `json:"function"`
}

type ToolCallFunction struct {
	Index     int                       `json:"index"`
	Name      string                    `json:"name"`
	Arguments ToolCallFunctionArguments `json:"arguments"`
}

type ToolCallFunctionArguments struct {
	M map[string]any
}

func (t ToolCallFunctionArguments) ToMap() map[string]any {
	if t.M == nil {
		return nil
	}
	return maps.Clone(t.M)
}

func (t *ToolCallFunctionArguments) UnmarshalJSON(data []byte) error {
	if t.M == nil {
		t.M = map[string]any{}
	}

	return json.Unmarshal(data, &t.M)
}

func (t ToolCallFunctionArguments) MarshalJSON() ([]byte, error) {
	if t.M == nil {
		return []byte("{}"), nil
	}

	return json.Marshal(t.M)
}

type Tool struct {
	Type     string       `json:"type"`
	Items    any          `json:"items,omitempty"`
	Function ToolFunction `json:"function"`
}

type PropertyType []string

func (pt *PropertyType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*pt = []string{s}
		return nil
	}

	var a []string
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*pt = a

	return nil
}

func (pt PropertyType) MarshalJSON() ([]byte, error) {
	if len(pt) == 1 {
		return json.Marshal(pt[0])
	}

	return json.Marshal([]string(pt))
}

func (pt PropertyType) String() string {
	if len(pt) == 0 {
		return ""
	}

	if len(pt) == 1 {
		return pt[0]
	}

	return fmt.Sprintf("%v", []string(pt))
}

type ToolPropertiesMap struct {
	m map[string]ToolProperty
}

func (t *ToolPropertiesMap) ToMap() map[string]ToolProperty {
	if t == nil || t.m == nil {
		return nil
	}

	return maps.Clone(t.m)
}

func (t *ToolPropertiesMap) UnmarshalJSON(data []byte) error {
	t.m = map[string]ToolProperty{}
	return json.Unmarshal(data, &t.m)
}

func (t ToolPropertiesMap) MarshalJSON() ([]byte, error) {
	if t.m == nil {
		return []byte("null"), nil
	}

	return json.Marshal(t.m)
}

type ToolProperty struct {
	AnyOf       []ToolProperty     `json:"anyOf,omitempty"`
	Type        PropertyType       `json:"type,omitempty"`
	Items       any                `json:"items,omitempty"`
	Description string             `json:"description,omitempty"`
	Enum        []any              `json:"enum,omitempty"`
	Properties  *ToolPropertiesMap `json:"properties,omitempty"`
}

func (tp ToolProperty) ToTypeScriptType() string {
	if len(tp.AnyOf) > 0 {
		var types []string
		for _, anyOf := range tp.AnyOf {
			types = append(types, anyOf.ToTypeScriptType())
		}

		return strings.Join(types, " | ")
	}

	if len(tp.Type) == 0 {
		return "any"
	}

	if len(tp.Type) == 1 {
		return mapToTypeScriptType(tp.Type[0])
	}

	var types []string
	for _, t := range tp.Type {
		types = append(types, mapToTypeScriptType(t))
	}

	return strings.Join(types, " | ")
}

func mapToTypeScriptType(jsonType string) string {
	switch jsonType {
	case "string":
		return "string"
	case "number", "integer":
		return "number"
	case "boolean":
		return "boolean"
	case "array":
		return "any[]"
	case "object":
		return "Record<string, any>"
	case "null":
		return "null"
	default:
		return "any"
	}
}

type ToolFunctionParameters struct {
	Type       string             `json:"type"`
	Defs       any                `json:"$defs,omitempty"`
	Items      any                `json:"items,omitempty"`
	Required   []string           `json:"required,omitempty"`
	Properties *ToolPropertiesMap `json:"properties"`
}

type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  ToolFunctionParameters `json:"parameters"`
}

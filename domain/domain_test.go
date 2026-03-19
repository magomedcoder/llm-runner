package domain

import "testing"

func TestFromProtoRole_ToProtoRole(t *testing.T) {
	tests := []struct {
		proto string
		want  AIChatMessageRole
	}{
		{"system", AIChatMessageRoleSystem},
		{"user", AIChatMessageRoleUser},
		{"assistant", AIChatMessageRoleAssistant},
		{"unknown", AIChatMessageRoleUser},
		{"", AIChatMessageRoleUser},
	}
	for _, tt := range tests {
		got := AIFromProtoRole(tt.proto)
		if got != tt.want {
			t.Errorf("FromProtoRole(%q) = %v, ожидалось %v", tt.proto, got, tt.want)
		}
		if back := AIToProtoRole(got); back != string(tt.want) && tt.proto != "unknown" && tt.proto != "" {
			if (tt.proto == "unknown" || tt.proto == "") && back == "user" {
				continue
			}
			t.Errorf("ToProtoRole(FromProtoRole(%q)) = %q", tt.proto, back)
		}
	}
}

func TestMessage_ToMap(t *testing.T) {
	m := &AIChatMessage{
		Content: "hi",
		Role:    AIChatMessageRoleUser,
	}
	out := m.AIToMap()
	if out["role"] != "user" || out["content"] != "hi" {
		t.Errorf("ToMap() вернул %v", out)
	}
}

func TestNewMessage(t *testing.T) {
	m := NewAIChatMessage(1, "text", AIChatMessageRoleUser)
	if m.SessionId != 1 || m.Content != "text" || m.Role != AIChatMessageRoleUser {
		t.Errorf("NewMessage: неверные поля %+v", m)
	}
}

func TestNewMessageWithAttachment(t *testing.T) {
	m := NewAIChatMessageWithAttachment(1, "text", AIChatMessageRoleAssistant, "file.txt", 0)
	if m.AttachmentName != "file.txt" || m.Role != AIChatMessageRoleAssistant {
		t.Errorf("NewMessageWithAttachment: неверные поля %+v", m)
	}
}

func TestNewChatSession(t *testing.T) {
	s := NewAIChatSession(1, "title", "model1")
	if s.UserId != 1 || s.Title != "title" || s.Model != "model1" {
		t.Errorf("NewChatSession: неверные поля %+v", s)
	}
}

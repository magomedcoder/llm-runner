package domain

import (
	"time"

	"github.com/magomedcoder/llm-runner/pb"
)

type AIChatMessageRole string

const (
	AIChatMessageRoleSystem    AIChatMessageRole = "system"
	AIChatMessageRoleUser      AIChatMessageRole = "user"
	AIChatMessageRoleAssistant AIChatMessageRole = "assistant"
)

type AIChatMessage struct {
	Id        int64
	SessionId int64
	Content   string
	Role      AIChatMessageRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

func AIFromProtoRole(role string) AIChatMessageRole {
	switch role {
	case "system":
		return AIChatMessageRoleSystem
	case "user":
		return AIChatMessageRoleUser
	case "assistant":
		return AIChatMessageRoleAssistant
	default:
		return AIChatMessageRoleUser
	}
}

func AIMessageFromProto(proto *pb.ChatMessage, sessionID int64) *AIChatMessage {
	if proto == nil {
		return nil
	}

	return &AIChatMessage{
		Id:        proto.Id,
		SessionId: sessionID,
		Content:   proto.Content,
		Role:      AIFromProtoRole(proto.Role),
		CreatedAt: time.Unix(proto.CreatedAt, 0),
		UpdatedAt: time.Unix(proto.CreatedAt, 0),
	}
}

func AIMessagesFromProto(protos []*pb.ChatMessage, sessionID int64) []*AIChatMessage {
	if len(protos) == 0 {
		return nil
	}

	out := make([]*AIChatMessage, len(protos))
	for i, p := range protos {
		out[i] = AIMessageFromProto(p, sessionID)
	}

	return out
}

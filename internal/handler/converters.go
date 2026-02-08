package handler

import (
	"github.com/magomedcoder/gen/api/pb"
	"github.com/magomedcoder/gen/internal/domain"
	"strconv"
)

func (a *AuthHandler) userToProto(user *domain.User) *pb.User {
	return &pb.User{
		Id:    strconv.Itoa(user.Id),
		Email: user.Email,
		Name:  user.Name,
	}
}

func (h *ChatHandler) sessionToProto(session *domain.ChatSession) *pb.ChatSession {
	return &pb.ChatSession{
		Id:        session.Id,
		Title:     session.Title,
		CreatedAt: session.CreatedAt.Unix(),
		UpdatedAt: session.UpdatedAt.Unix(),
	}
}

func (h *ChatHandler) messageToProto(msg *domain.Message) *pb.ChatMessage {
	return &pb.ChatMessage{
		Id:        msg.Id,
		Content:   msg.Content,
		Role:      domain.ToProtoRole(msg.Role),
		CreatedAt: msg.CreatedAt.Unix(),
	}
}

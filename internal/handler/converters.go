package handler

import (
	"strconv"

	"github.com/magomedcoder/gen/api/pb"
	"github.com/magomedcoder/gen/internal/domain"
)

func userToProto(user *domain.User) *pb.User {
	return &pb.User{
		Id:       strconv.Itoa(user.Id),
		Username: user.Username,
		Name:     user.Name,
		Surname:  user.Surname,
		Role:     int32(user.Role),
	}
}

func (c *ChatHandler) sessionToProto(session *domain.ChatSession) *pb.ChatSession {
	return &pb.ChatSession{
		Id:        session.Id,
		Title:     session.Title,
		CreatedAt: session.CreatedAt.Unix(),
		UpdatedAt: session.UpdatedAt.Unix(),
	}
}

func (c *ChatHandler) messageToProto(msg *domain.Message) *pb.ChatMessage {
	return &pb.ChatMessage{
		Id:        msg.Id,
		Content:   msg.Content,
		Role:      domain.ToProtoRole(msg.Role),
		CreatedAt: msg.CreatedAt.Unix(),
	}
}

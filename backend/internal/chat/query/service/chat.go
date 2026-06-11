package service

import (
	"backend/internal/chat/readmodel"
	"context"

	"github.com/google/uuid"
)

type ChatQueryRepository interface {
	ListMessagesByRoom(ctx context.Context, roomID uuid.UUID, limit int) ([]*readmodel.ChatReadModel, error)
}

type ChatService struct {
	repository ChatQueryRepository
}

func NewChatService(repository ChatQueryRepository) *ChatService {
	return &ChatService{
		repository: repository,
	}
}

func (s *ChatService) ListMessage(ctx context.Context, room string) ([]*readmodel.ChatReadModel, error) {
	roomID, err := uuid.Parse(room)
	if err != nil {
		return nil, err
	}
	result, err := s.repository.ListMessagesByRoom(ctx, roomID, 100)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return make([]*readmodel.ChatReadModel, 0), nil
	}
	return result, nil
}

// GetHistory は編集履歴を返します。ScyllaDB への移行は未対応のため空を返します。
func (s *ChatService) GetHistory(ctx context.Context, room string, chatId string) ([]*readmodel.ChatReadModel, error) {
	return make([]*readmodel.ChatReadModel, 0), nil
}

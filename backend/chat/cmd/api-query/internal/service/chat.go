package service

import (
	"boardroom/chat-readmodel"
	"context"

	"github.com/google/uuid"
)

type ChatQueryRepository interface {
	ListMessagesByRoom(ctx context.Context, roomID uuid.UUID, limit int) ([]*readmodel.Chat, error)
	ListMessageHistories(ctx context.Context, roomID uuid.UUID, chatID uuid.UUID) ([]*readmodel.Chat, error)
}

type ChatService struct {
	repository ChatQueryRepository
}

func NewChatService(repository ChatQueryRepository) *ChatService {
	return &ChatService{
		repository: repository,
	}
}

func (s *ChatService) ListMessage(ctx context.Context, roomID uuid.UUID) ([]*readmodel.Chat, error) {
	// ToDo validation room.visible(user)
	result, err := s.repository.ListMessagesByRoom(ctx, roomID, 100)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return make([]*readmodel.Chat, 0), nil
	}
	return result, nil
}

// GetHistory は編集履歴を返します。ScyllaDB への移行は未対応のため空を返します。
func (s *ChatService) GetHistory(ctx context.Context, roomID uuid.UUID, chatID uuid.UUID) ([]*readmodel.Chat, error) {
	// ToDo validation room.visible(user)
	return s.repository.ListMessageHistories(ctx, roomID, chatID)
}

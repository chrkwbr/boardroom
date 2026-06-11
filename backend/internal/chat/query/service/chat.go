package service

import (
	"backend/internal/chat/query/repository"
	"backend/internal/chat/readmodel"
	"context"
	"slices"
)

type ChatService struct {
	repository *repository.ChatReadModelRepository
}

func NewChatService(repository *repository.ChatReadModelRepository) *ChatService {
	return &ChatService{
		repository: repository,
	}
}

func (s *ChatService) ListMessage(ctx context.Context, room string) ([]*readmodel.ChatReadModel, error) {
	chatIds, err := s.repository.ZRevRangeRoomChatIds(ctx, room, 0, 99)
	if err != nil {
		return nil, err
	}
	if len(chatIds) == 0 {
		return make([]*readmodel.ChatReadModel, 0), nil
	}
	slices.Reverse(chatIds)

	return s.repository.MGetChat(ctx, chatIds)
}

func (s *ChatService) GetHistory(ctx context.Context, room string, chatId string) ([]*readmodel.ChatReadModel, error) {
	return s.repository.LRangeHistory(ctx, chatId)
}

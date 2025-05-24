package service

import (
	"backend/chat/query"
	"backend/chat/query/repository"
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

func (s *ChatService) ListMessage(ctx context.Context, room string) ([]*query.ChatReadModel, error) {
	chatIds, err := s.repository.ZRevRangeRoomChatIds(ctx, room, 0, 99)
	if err != nil {
		return nil, err
	}
	if len(chatIds) == 0 {
		return make([]*query.ChatReadModel, 0), nil
	}
	slices.Reverse(chatIds)

	return s.repository.MGetChat(ctx, chatIds)
}

func (s *ChatService) GetHistory(ctx context.Context, room string, chatId string) ([]*query.ChatReadModel, error) {
	return s.repository.LRangeHistory(ctx, chatId)
}

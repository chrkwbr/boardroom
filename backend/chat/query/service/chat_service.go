package service

import (
	"backend/chat/command/domain"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

type ChatService struct {
	redisClient *redis.Client
}

func NewChatService(redisClient *redis.Client) *ChatService {
	return &ChatService{
		redisClient: redisClient,
	}
}

func (s *ChatService) ListMessage(ctx context.Context, room string) ([]domain.Chat, error) {
	key := "chat:room:myroom"
	//key := "chat:room:" + room

	result, err := s.redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var messages []domain.Chat
	for _, msg := range result {
		chat := &domain.Chat{}
		if err := json.Unmarshal([]byte(msg), chat); err != nil {
			return nil, err
		}
		messages = append(messages, *chat)
	}
	return messages, nil
}

package service

import (
	"backend/chat/command/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"slices"
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
	key := fmt.Sprintf("chats:%v", room)
	chatIds, err := s.redisClient.ZRevRange(ctx, key, 0, 99).Result()
	if err != nil {
		return nil, err
	}
	if len(chatIds) == 0 {
		return make([]domain.Chat, 0), nil
	}
	slices.Reverse(chatIds)

	chatKeys := make([]string, len(chatIds))
	for i, chatId := range chatIds {
		chatKeys[i] = fmt.Sprintf("chat:%s", chatId)
	}

	chats, err := s.redisClient.MGet(ctx, chatKeys...).Result()
	if err != nil {
		return nil, err
	}

	var messages []domain.Chat
	for _, c := range chats {
		chat := &domain.Chat{}
		if err := json.Unmarshal([]byte(c.(string)), chat); err != nil {
			log.Println("Failed to unmarshal chat:", err)
			continue
		}
		messages = append(messages, *chat)
	}
	return messages, nil
}

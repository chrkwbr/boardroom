package repository

import (
	"boardroom/chat-readmodel"
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	formatChatKey    = "chat:message:%s"
	formatHistoryKey = "chat:room:%v:timeline"
	formatChatIDsKey = "chat:message:%v:history"
)

type ChatReadModelRepository struct {
	redis *redis.Client
}

func NewChatReadModelRepository(redisClient *redis.Client) *ChatReadModelRepository {
	return &ChatReadModelRepository{
		redis: redisClient,
	}
}

// Chat
func (r *ChatReadModelRepository) MGetChat(ctx context.Context, chatIds []string) ([]*readmodel.ChatReadModel, error) {
	keys := make([]string, len(chatIds))
	for i, chatId := range chatIds {
		keys[i] = fmt.Sprintf(formatChatKey, chatId)
	}

	results, err := r.redis.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	var readModels []*readmodel.ChatReadModel
	for _, result := range results {
		if result == nil || result == "" {
			continue
		}
		readModel := &readmodel.ChatReadModel{}
		if err := json.Unmarshal([]byte(result.(string)), readModel); err != nil {
			return nil, err
		}
		readModels = append(readModels, readModel)
	}
	return readModels, nil
}

func (r *ChatReadModelRepository) ZRevRangeRoomChatIds(ctx context.Context, roomId string, start, end int64) ([]string, error) {
	chatRoomKey := fmt.Sprintf(formatChatIDsKey, roomId)
	chatIds, err := r.redis.ZRevRange(ctx, chatRoomKey, start, end).Result()
	if err != nil {
		return nil, err
	}
	return chatIds, nil
}

// History

func (r *ChatReadModelRepository) LRangeHistory(ctx context.Context, chatId string) ([]*readmodel.ChatReadModel, error) {
	chatHistoryKey := fmt.Sprintf(formatHistoryKey, chatId)
	result, err := r.redis.LRange(ctx, chatHistoryKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var messages []*readmodel.ChatReadModel
	for _, c := range result {
		chat := &readmodel.ChatReadModel{}
		if err := json.Unmarshal([]byte(c), chat); err != nil {
			return nil, err
		}
		messages = append(messages, chat)
	}
	return messages, nil
}

package repository

import (
	"backend/chat/event"
	"backend/chat/query"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type ChatReadModelRepository struct {
	redis *redis.Client
}

func NewChatReadModelRepository(redisClient *redis.Client) *ChatReadModelRepository {
	return &ChatReadModelRepository{
		redis: redisClient,
	}
}

const (
	formatChatKey    = "chat:%s"
	formatHistoryKey = "chats:%v:history" // chats:${chatId}:history
	formatCHatIdsKey = "chats:%v"         // chats:${roomId}
)

// Chat

func (r *ChatReadModelRepository) SetChat(ctx context.Context, chatEvent *event.ChatEvent) (*query.ChatReadModel, error) {
	key := fmt.Sprintf(formatChatKey, chatEvent.ChatId)
	readModel, err := query.FromPayload(chatEvent)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(readModel)
	if err != nil {
		return nil, err
	}

	if err := r.redis.Set(ctx, key, payload, time.Hour*24*10).Err(); err != nil {
		return nil, err
	}
	return readModel, nil
}

func (r *ChatReadModelRepository) SetArgsChat(ctx context.Context, chatEvent *event.ChatEvent) (*query.ChatReadModel, error) {
	key := fmt.Sprintf(formatChatKey, chatEvent.ChatId)
	readModel, err := query.FromPayload(chatEvent)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(readModel)
	if err != nil {
		return nil, err
	}

	previewChat, err := r.redis.SetArgs(ctx, key, payload, redis.SetArgs{
		Get: true,
		TTL: time.Hour * 24 * 10,
	}).Result()
	if err != nil {
		return nil, err
	}

	preview := &query.ChatReadModel{}
	if err := json.Unmarshal([]byte(previewChat), &preview); err != nil {
		return nil, err
	}

	return preview, nil
}

func (r *ChatReadModelRepository) MGetChat(ctx context.Context, chatIds []string) ([]*query.ChatReadModel, error) {
	keys := make([]string, len(chatIds))
	for i, chatId := range chatIds {
		keys[i] = fmt.Sprintf(formatChatKey, chatId)
	}

	results, err := r.redis.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	var readModels []*query.ChatReadModel
	for _, result := range results {
		if result == nil || result == "" {
			continue
		}
		readModel := &query.ChatReadModel{}
		if err := json.Unmarshal([]byte(result.(string)), readModel); err != nil {
			return nil, err
		}
		readModels = append(readModels, readModel)
	}
	return readModels, nil
}

func (r *ChatReadModelRepository) DelChat(ctx context.Context, chatId uuid.UUID) error {
	key := fmt.Sprintf(formatChatKey, chatId)
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

// Room Chat IDs

func (r *ChatReadModelRepository) ZAddNXRoomChatIds(ctx context.Context, readModel *query.ChatReadModel) error {
	chatRoomKey := fmt.Sprintf(formatCHatIdsKey, readModel.Room)
	z := redis.Z{
		Score:  float64(readModel.CreatedAt),
		Member: fmt.Sprintf("%s", readModel.ID),
	}
	if err := r.redis.ZAddNX(ctx, chatRoomKey, z).Err(); err != nil {
		return err
	}
	return nil
}

func (r *ChatReadModelRepository) ZRevRangeRoomChatIds(ctx context.Context, roomId string, start, end int64) ([]string, error) {
	chatRoomKey := fmt.Sprintf(formatCHatIdsKey, roomId)
	chatIds, err := r.redis.ZRevRange(ctx, chatRoomKey, start, end).Result()
	if err != nil {
		return nil, err
	}
	return chatIds, nil
}

func (r *ChatReadModelRepository) ZRemRoomChatIds(ctx context.Context, roomId string, chatId uuid.UUID) error {
	chatRoomKey := fmt.Sprintf(formatCHatIdsKey, roomId)
	if err := r.redis.ZRem(ctx, chatRoomKey, fmt.Sprintf("%s", chatId)).Err(); err != nil {
		return err
	}
	return nil
}

// History

func (r *ChatReadModelRepository) LPushHistory(ctx context.Context, readModel *query.ChatReadModel) error {
	m, err := json.Marshal(readModel)
	if err != nil {
		return err
	}
	r.redis.LPush(ctx, fmt.Sprintf(formatHistoryKey, readModel.ID), m)
	return nil
}

func (r *ChatReadModelRepository) LRangeHistory(ctx context.Context, chatId string) ([]*query.ChatReadModel, error) {
	chatHistoryKey := fmt.Sprintf(formatHistoryKey, chatId)
	result, err := r.redis.LRange(ctx, chatHistoryKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var messages []*query.ChatReadModel
	for _, c := range result {
		chat := &query.ChatReadModel{}
		if err := json.Unmarshal([]byte(c), chat); err != nil {
			return nil, err
		}
		messages = append(messages, chat)
	}
	return messages, nil
}

func (r *ChatReadModelRepository) DelHistory(ctx context.Context, chatId uuid.UUID) error {
	chatHistoryKey := fmt.Sprintf(formatHistoryKey, chatId)
	if err := r.redis.Del(ctx, chatHistoryKey).Err(); err != nil {
		return err
	}
	return nil
}

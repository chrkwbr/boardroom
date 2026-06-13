package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ChatRedisRepository struct {
	redis *redis.Client
}

func NewChatRedisRepository(redisClient *redis.Client) *ChatRedisRepository {
	return &ChatRedisRepository{
		redis: redisClient,
	}
}

const (
	FormatChatKey     = "chat:message:%s"
	FormatHistoryKey  = "chat:room:%v:timeline"
	FormatCHatIdsKey  = "chat:message:%v:history"
	FormatChatChannel = "chat:room:%s:updates"
)

// Chat
func (r *ChatRedisRepository) GetChat(ctx context.Context, chatId uuid.UUID) (*Chat, error) {
	result, err := r.redis.Get(ctx, fmt.Sprintf(FormatChatKey, chatId)).Result()
	if err != nil {
		return nil, err
	}
	readModel := &Chat{}
	if err := json.Unmarshal([]byte(result), readModel); err != nil {
		return nil, err
	}
	return readModel, nil
}

func (r *ChatRedisRepository) SetChat(ctx context.Context, model *Chat) error {
	key := fmt.Sprintf(FormatChatKey, model.ID)
	payload, err := json.Marshal(model)
	if err != nil {
		return err
	}

	if err := r.redis.Set(ctx, key, payload, time.Hour*24*10).Err(); err != nil {
		return err
	}
	return nil
}

func (r *ChatRedisRepository) PublishChatEvent(ctx context.Context, roomId uuid.UUID, event *ChatRedisEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	channel := fmt.Sprintf(FormatChatChannel, roomId)
	if err := r.redis.Publish(ctx, channel, payload).Err(); err != nil {
		return err
	}
	return nil
}

//	func (r *ChatRedisRepository) SetArgsChat(ctx context.Context, chatEvent *domain.ChatEvent) (*query.Chat, error) {
//		key := fmt.Sprintf(FormatChatKey, chatEvent.ID)
//		readModel, err := query.FromPayload(chatEvent)
//		if err != nil {
//			return nil, err
//		}
//		payload, err := json.Marshal(readModel)
//		if err != nil {
//			return nil, err
//		}
//
//		previewChat, err := r.redis.SetArgs(ctx, key, payload, redis.SetArgs{
//			Get: true,
//			TTL: time.Hour * 24 * 10,
//		}).Result()
//		if err != nil {
//			return nil, err
//		}
//
//		preview := &query.Chat{}
//		if err := json.Unmarshal([]byte(previewChat), &preview); err != nil {
//			return nil, err
//		}
//
//		return preview, nil
//	}

func (r *ChatRedisRepository) DelChat(ctx context.Context, chatId uuid.UUID) error {
	key := fmt.Sprintf(FormatChatKey, chatId)
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

// Room Chat IDs
func (r *ChatRedisRepository) ZAddNXRoomChatIds(ctx context.Context, model *Chat) error {
	chatRoomKey := fmt.Sprintf(FormatCHatIdsKey, model.RoomID)
	z := redis.Z{
		Score:  float64(model.CreatedAt),
		Member: fmt.Sprintf("%s", model.ID),
	}
	if err := r.redis.ZAddNX(ctx, chatRoomKey, z).Err(); err != nil {
		return err
	}
	return nil
}

func (r *ChatRedisRepository) ZRemRoomChatIds(ctx context.Context, roomId uuid.UUID, chatId uuid.UUID) error {
	chatRoomKey := fmt.Sprintf(FormatCHatIdsKey, roomId)
	if err := r.redis.ZRem(ctx, chatRoomKey, fmt.Sprintf("%s", chatId)).Err(); err != nil {
		return err
	}
	return nil
}

// History
func (r *ChatRedisRepository) LPushHistory(ctx context.Context, model *Chat) error {
	m, err := json.Marshal(model)
	if err != nil {
		return err
	}
	r.redis.LPush(ctx, fmt.Sprintf(FormatHistoryKey, model.ID), m)
	return nil
}

func (r *ChatRedisRepository) DelHistory(ctx context.Context, chatId uuid.UUID) error {
	chatHistoryKey := fmt.Sprintf(FormatHistoryKey, chatId)
	if err := r.redis.Del(ctx, chatHistoryKey).Err(); err != nil {
		return err
	}
	return nil
}

package internal

import (
	"boardroom/chat-notification"
	"boardroom/shared/event"
	"boardroom/shared/infra/pubsub"
	"context"
	"encoding/json"
	"log"
	"time"
)

type ChatNotificationPublisher struct {
	subscriber pubsub.EventSubscriber
	repo       *notification.ChatRedisRepository
}

func NewChatNotificationPublisher(sub pubsub.EventSubscriber, repo *notification.ChatRedisRepository) *ChatNotificationPublisher {
	return &ChatNotificationPublisher{
		subscriber: sub,
		repo:       repo,
	}
}

func (rc *ChatNotificationPublisher) Start() {
	go func() {
		const maxRetries = 10
		for i := range maxRetries {
			err := rc.subscriber.Subscribe("_", func(key string, value []byte) error {
				rc.process(context.Background(), value)
				return nil
			})
			if err == nil {
				return
			}
			wait := time.Duration(i+1) * 3 * time.Second
			log.Printf("Failed to subscribe (attempt %d/%d): %v — retrying in %s", i+1, maxRetries, err, wait)
			time.Sleep(wait)
		}
		log.Println("ChatNotificationPublisher: gave up after max retries")
	}()
}

func (rc *ChatNotificationPublisher) process(ctx context.Context, msg []byte) {
	chatEvent := &event.ChatEvent{}
	if err := json.Unmarshal(msg, chatEvent); err != nil {
		log.Println("Failed to unmarshal chat event:", err)
		return
	}

	switch chatEvent.Type {
	case event.EventTypeCreated:
		rc.onCreate(ctx, chatEvent)
	case event.EventTypeUpdated:
		rc.onUpdate(ctx, chatEvent)
	case event.EventTypeDeleted:
		rc.onDelete(ctx, chatEvent)
	default:
		log.Println("Unknown event type:", chatEvent.Type)
	}
}

func (rc *ChatNotificationPublisher) onCreate(ctx context.Context, evt *event.ChatEvent) {
	p := &event.ChatCreatedPayload{}
	if err := json.Unmarshal(evt.Payload, p); err != nil {
		log.Println("Failed to unmarshal ChatCreatedPayload:", err)
		return
	}

	// ToDo: sender 情報は将来的に Redis から取得
	model := &notification.Chat{
		ID: p.ID,
		Sender: notification.User{
			ID:   p.SenderID,
			Name: "test name",
			Icon: "https://img.daisyui.com/images/profile/demo/1@94.webp",
		},
		RoomID:    p.RoomID,
		Message:   p.Message,
		Version:   p.Version,
		CreatedAt: evt.OccurredAt,
		UpdatedAt: evt.OccurredAt,
	}

	if err := rc.repo.SetChat(ctx, model); err != nil {
		log.Println("Failed to save chat to Redis:", err)
		return
	}

	e := notification.NewChatCreatedEvent(*model)
	if err := rc.repo.PublishChatEvent(ctx, model.RoomID, e); err != nil {
		log.Println("Failed to publish chat event:", err)
	}

	//if err := rc.repo.LPushHistory(ctx, model); err != nil {
	//	log.Println("Failed to push chat history:", err)
	//	return
	//}
	//
	//if err := rc.repo.ZAddNXRoomChatIds(ctx, model); err != nil {
	//	log.Println("Failed to save room chat IDs:", err)
	//	return
	//}
}

func (rc *ChatNotificationPublisher) onUpdate(ctx context.Context, evt *event.ChatEvent) {
	p := &event.ChatEditedPayload{}
	if err := json.Unmarshal(evt.Payload, p); err != nil {
		log.Println("Failed to unmarshal ChatEditedPayload:", err)
		return
	}

	orig, err := rc.repo.GetChat(ctx, p.ID)
	if err != nil {
		log.Println("Failed to get chat from Redis:", err)
		return
	}

	edited := orig.NewUpdate(p.Message, evt.OccurredAt)

	if err := rc.repo.SetChat(ctx, edited); err != nil {
		log.Println("Failed to update chat in Redis:", err)
		return
	}

	e := notification.NewChatEditedEvent(*edited)
	if err := rc.repo.PublishChatEvent(ctx, edited.RoomID, e); err != nil {
		log.Println("Failed to publish chat event:", err)
	}

	//if err := rc.repo.LPushHistory(ctx, edited); err != nil {
	//	log.Println("Failed to push chat history:", err)
	//	return
	//}
	//
	//if err := rc.repo.ZAddNXRoomChatIds(ctx, edited); err != nil {
	//	log.Println("Failed to save room chat IDs:", err)
	//	return
	//}
}

func (rc *ChatNotificationPublisher) onDelete(ctx context.Context, evt *event.ChatEvent) {
	p := &event.ChatDeletedPayload{}
	if err := json.Unmarshal(evt.Payload, p); err != nil {
		log.Println("Failed to unmarshal ChatDeletedPayload:", err)
		return
	}

	orig, err := rc.repo.GetChat(ctx, p.ID)
	if err != nil {
		log.Println("Failed to get chat from Redis:", err)
		return
	}

	e := notification.NewChatDeletedEvent(p.RoomID, p.ID)
	if err := rc.repo.PublishChatEvent(ctx, orig.RoomID, e); err != nil {
		log.Println("Failed to publish chat event:", err)
	}

	if err := rc.repo.DelChat(ctx, p.ID); err != nil {
		log.Println("Failed to delete chat from Redis:", err)
	}

	//if err := rc.repo.ZRemRoomChatIds(ctx, orig.RoomID, p.ID); err != nil {
	//	log.Println("Failed to remove chat ID from room:", err)
	//	return
	//}
	//
	//if err := rc.repo.DelHistory(ctx, p.ID); err != nil {
	//	log.Println("Failed to delete chat history:", err)
	//	return
	//}

}

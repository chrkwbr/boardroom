package readmodel

import (
	"backend/internal/chat/domain"
	"backend/internal/shared/infra/pubsub"
	"context"
	"encoding/json"
	"log"
)

type Materializer struct {
	subscriber pubsub.EventSubscriber
	scylla     *ChatScyllaRepository
}

func NewMaterializer(sub pubsub.EventSubscriber, scylla *ChatScyllaRepository) *Materializer {
	return &Materializer{
		subscriber: sub,
		scylla:     scylla,
	}
}

func (m *Materializer) Start() {
	go func() {
		if err := m.subscriber.Subscribe("_", func(key string, value []byte) error {
			m.process(context.Background(), value)
			return nil
		}); err != nil {
			log.Panicln("Failed to subscribe to event:", err)
		}
	}()
}

func (m *Materializer) process(ctx context.Context, msg []byte) {
	chatEvent := &domain.ChatEvent{}
	if err := json.Unmarshal(msg, chatEvent); err != nil {
		log.Println("Failed to unmarshal chat event:", err)
		return
	}

	switch chatEvent.Type {
	case domain.EventTypeCreated:
		m.onCreate(ctx, chatEvent)
	case domain.EventTypeUpdated:
		m.onUpdate(ctx, chatEvent)
	case domain.EventTypeDeleted:
		m.onDelete(ctx, chatEvent)
	default:
		log.Println("Unknown event type:", chatEvent.Type)
	}
}

func (m *Materializer) onCreate(ctx context.Context, event *domain.ChatEvent) {
	p := &domain.ChatCreatedPayload{}
	if err := json.Unmarshal(event.Payload, p); err != nil {
		log.Println("Failed to unmarshal ChatCreatedPayload:", err)
		return
	}

	// ToDo: sender 情報は将来的に User サービスから取得
	model := &ChatReadModel{
		ID: p.ID,
		Sender: User{
			ID:   p.SenderID,
			Name: "test name",
			Icon: "https://img.daisyui.com/images/profile/demo/1@94.webp",
		},
		RoomID:    p.RoomID,
		Message:   p.Message,
		Version:   p.Version,
		CreatedAt: event.OccurredAt,
		UpdatedAt: event.OccurredAt,
	}

	if err := m.scylla.InsertChat(ctx, model); err != nil {
		log.Println("Failed to insert chat to ScyllaDB:", err)
		return
	}
}

func (m *Materializer) onUpdate(ctx context.Context, event *domain.ChatEvent) {
	p := &domain.ChatEditedPayload{}
	if err := json.Unmarshal(event.Payload, p); err != nil {
		log.Println("Failed to unmarshal ChatEditedPayload:", err)
		return
	}

	orig, err := m.scylla.GetChatByID(ctx, p.ID)
	if err != nil {
		log.Println("Failed to get original chat from ScyllaDB:", err)
		return
	}

	edited, _ := orig.NewUpdate(p.Message, event.OccurredAt)

	if err := m.scylla.UpdateChat(ctx, edited); err != nil {
		log.Println("Failed to update chat in ScyllaDB:", err)
		return
	}
}

func (m *Materializer) onDelete(ctx context.Context, event *domain.ChatEvent) {
	p := &domain.ChatDeletedPayload{}
	if err := json.Unmarshal(event.Payload, p); err != nil {
		log.Println("Failed to unmarshal ChatDeletedPayload:", err)
		return
	}

	orig, err := m.scylla.GetChatByID(ctx, p.ID)
	if err != nil {
		log.Println("Failed to get original chat from ScyllaDB:", err)
		return
	}

	if err := m.scylla.DeleteChat(ctx, orig.ID); err != nil {
		log.Println("Failed to delete chat from ScyllaDB:", err)
	}
}

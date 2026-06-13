package internal

import (
	"boardroom/chat-readmodel"
	"boardroom/shared/event"
	"boardroom/shared/infra/pubsub"
	"context"
	"encoding/json"
	"log"
	"time"
)

type Materializer struct {
	subscriber pubsub.EventSubscriber
	scylla     *readmodel.ChatScyllaRepository
}

func NewMaterializer(sub pubsub.EventSubscriber, scylla *readmodel.ChatScyllaRepository) *Materializer {
	return &Materializer{
		subscriber: sub,
		scylla:     scylla,
	}
}

func (m *Materializer) Start() {
	go func() {
		const maxRetries = 10
		for i := range maxRetries {
			err := m.subscriber.Subscribe("_", func(key string, value []byte) error {
				m.process(context.Background(), value)
				return nil
			})
			if err == nil {
				return
			}
			wait := time.Duration(i+1) * 3 * time.Second
			log.Printf("Materializer: subscribe failed (attempt %d/%d): %v — retrying in %s", i+1, maxRetries, err, wait)
			time.Sleep(wait)
		}
	}()
}

func (m *Materializer) process(ctx context.Context, msg []byte) {
	chatEvent := &event.ChatEvent{}
	if err := json.Unmarshal(msg, chatEvent); err != nil {
		log.Println("Failed to unmarshal chat event:", err)
		return
	}

	switch chatEvent.Type {
	case event.EventTypeCreated:
		m.onCreate(ctx, chatEvent)
	case event.EventTypeUpdated:
		m.onUpdate(ctx, chatEvent)
	case event.EventTypeDeleted:
		m.onDelete(ctx, chatEvent)
	default:
		log.Println("Unknown event type:", chatEvent.Type)
	}
}

func (m *Materializer) onCreate(ctx context.Context, evt *event.ChatEvent) {
	p := &event.ChatCreatedPayload{}
	if err := json.Unmarshal(evt.Payload, p); err != nil {
		log.Println("Failed to unmarshal ChatCreatedPayload:", err)
		return
	}

	model := &readmodel.Chat{
		ID:        p.ID,
		SenderID:  p.SenderID,
		RoomID:    p.RoomID,
		Message:   p.Message,
		Version:   p.Version,
		CreatedAt: evt.OccurredAt,
		UpdatedAt: evt.OccurredAt,
	}

	if err := m.scylla.InsertChat(ctx, model); err != nil {
		log.Println("Failed to insert chat to ScyllaDB:", err)
		return
	}

	if err := m.scylla.InsertHistory(ctx, model, readmodel.Created); err != nil {
		log.Println("Failed to insert history to ScyllaDB:", err)
		return
	}
}

func (m *Materializer) onUpdate(ctx context.Context, evt *event.ChatEvent) {
	p := &event.ChatEditedPayload{}
	if err := json.Unmarshal(evt.Payload, p); err != nil {
		log.Println("Failed to unmarshal ChatEditedPayload:", err)
		return
	}

	orig, err := m.scylla.GetChat(ctx, p.RoomID, p.ID)
	if err != nil {
		log.Println("Failed to get original chat from ScyllaDB:", err)
		return
	}

	edited := orig.NewUpdate(p.Message, evt.OccurredAt)

	if err := m.scylla.UpdateChat(ctx, edited); err != nil {
		log.Println("Failed to update chat in ScyllaDB:", err)
		return
	}

	if err := m.scylla.InsertHistory(ctx, edited, readmodel.Edited); err != nil {
		log.Println("Failed to insert history in ScyllaDB:", err)
		return
	}
}

func (m *Materializer) onDelete(ctx context.Context, evt *event.ChatEvent) {
	p := &event.ChatDeletedPayload{}
	if err := json.Unmarshal(evt.Payload, p); err != nil {
		log.Println("Failed to unmarshal ChatDeletedPayload:", err)
		return
	}

	orig, err := m.scylla.GetChat(ctx, p.RoomID, p.ID)
	if err != nil {
		log.Println("Failed to get original chat from ScyllaDB:", err)
		return
	}
	del := orig.NewDelete()

	if err := m.scylla.DeleteChat(ctx, p.RoomID, p.ID); err != nil {
		log.Println("Failed to delete chat from ScyllaDB:", err)
	}

	if err = m.scylla.InsertHistory(ctx, del, readmodel.Deleted); err != nil {
		log.Println("Failed to insert history in ScyllaDB:", err)
		return
	}
}

package service

import (
	"boardroom/chat-readmodel"
	"context"

	"github.com/google/uuid"
)

type ChatQueryRepository interface {
	ListMessagesByRoom(ctx context.Context, roomID uuid.UUID, limit int) ([]*readmodel.Chat, error)
	ListMessageHistories(ctx context.Context, roomID uuid.UUID, chatID uuid.UUID) ([]*readmodel.Chat, error)
}

type ChatService struct {
	repository ChatQueryRepository
}

func NewChatService(repository ChatQueryRepository) *ChatService {
	return &ChatService{
		repository: repository,
	}
}

type Chat struct {
	RoomID    uuid.UUID `json:"roomId"`
	ID        uuid.UUID `json:"id"`
	Message   string    `json:"message"`
	Sender    User      `json:"sender"`
	Version   int64     `json:"version"`
	CreatedAt int64     `json:"createdAt"`
	UpdatedAt int64     `json:"updatedAt"`
}

type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Icon string    `json:"icon"`
}

func (s *ChatService) ListMessage(ctx context.Context, roomID uuid.UUID) ([]*Chat, error) {
	// ToDo validation room.visible(user)
	r, err := s.repository.ListMessagesByRoom(ctx, roomID, 100)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return make([]*Chat, 0), nil
	}

	results := make([]*Chat, len(r))
	user := make(map[uuid.UUID]*User, len(r))
	for i := range r {
		s := r[i].SenderID
		if _, ok := user[s]; !ok {
			// fetch from redis
			user[s] = &User{
				ID:   s,
				Name: "user",
				Icon: "https://img.daisyui.com/images/profile/demo/1@94.webp",
			}
		}
		results[i] = &Chat{
			RoomID:    r[i].RoomID,
			ID:        r[i].ID,
			Message:   r[i].Message,
			Sender:    *user[s],
			Version:   r[i].Version,
			CreatedAt: r[i].CreatedAt,
			UpdatedAt: r[i].UpdatedAt,
		}
	}
	return results, nil
}

// GetHistory は編集履歴を返します。ScyllaDB への移行は未対応のため空を返します。
func (s *ChatService) GetHistory(ctx context.Context, roomID uuid.UUID, chatID uuid.UUID) ([]*readmodel.Chat, error) {
	// ToDo validation room.visible(user)
	return s.repository.ListMessageHistories(ctx, roomID, chatID)
}

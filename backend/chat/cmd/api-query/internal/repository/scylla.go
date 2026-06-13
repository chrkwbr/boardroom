package repository

import (
	"boardroom/chat-readmodel"
	"context"

	"github.com/google/uuid"
)

type ChatScyllaQueryRepository struct {
	scylla *readmodel.ChatScyllaRepository
}

func NewChatScyllaQueryRepository(scylla *readmodel.ChatScyllaRepository) *ChatScyllaQueryRepository {
	return &ChatScyllaQueryRepository{scylla: scylla}
}

func (r *ChatScyllaQueryRepository) ListMessagesByRoom(ctx context.Context, roomID uuid.UUID, limit int) ([]*readmodel.Chat, error) {
	return r.scylla.GetChatsByRoomID(ctx, roomID, limit)
}

func (r *ChatScyllaQueryRepository) ListMessageHistories(ctx context.Context, roomID uuid.UUID, chatID uuid.UUID) ([]*readmodel.Chat, error) {
	return r.scylla.GetHistory(ctx, roomID, chatID)
}

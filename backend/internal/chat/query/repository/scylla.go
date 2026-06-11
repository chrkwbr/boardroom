package repository

import (
	"backend/internal/chat/readmodel"
	"context"

	"github.com/google/uuid"
)

type ChatScyllaQueryRepository struct {
	scylla *readmodel.ChatScyllaRepository
}

func NewChatScyllaQueryRepository(scylla *readmodel.ChatScyllaRepository) *ChatScyllaQueryRepository {
	return &ChatScyllaQueryRepository{scylla: scylla}
}

func (r *ChatScyllaQueryRepository) ListMessagesByRoom(ctx context.Context, roomID uuid.UUID, limit int) ([]*readmodel.ChatReadModel, error) {
	return r.scylla.GetChatsByRoomID(ctx, roomID, limit)
}


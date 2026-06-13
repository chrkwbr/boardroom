package repository

import (
	"backend/chat/pkg/shared/readmodel"
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

func (r *ChatScyllaQueryRepository) ListMessageHistories(ctx context.Context, roomID uuid.UUID, chatID uuid.UUID) ([]*readmodel.ChatReadModel, error) {
	return r.scylla.GetHistory(ctx, roomID, chatID)
}

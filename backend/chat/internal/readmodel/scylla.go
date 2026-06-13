package readmodel

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type ChatScyllaRepository struct {
	session *gocql.Session
}

func NewChatScyllaRepository(hosts ...string) (*ChatScyllaRepository, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = "chat"
	cluster.Consistency = gocql.One
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 10 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create scylla session: %w", err)
	}
	return &ChatScyllaRepository{session: session}, nil
}

func (r *ChatScyllaRepository) Close() {
	r.session.Close()
}

func toGocql(u uuid.UUID) gocql.UUID   { return gocql.UUID(u) }
func fromGocql(u gocql.UUID) uuid.UUID { return uuid.UUID(u) }

func (r *ChatScyllaRepository) InsertChat(ctx context.Context, m *Chat) error {
	return r.session.Query(`
		INSERT INTO chat.chat_messages
			(room_id, id, sender_id, message, version, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		toGocql(m.RoomID), toGocql(m.ID), toGocql(m.SenderID),
		m.Message, m.Version, m.CreatedAt, m.UpdatedAt,
	).WithContext(ctx).Exec()
}

// GetChat は room_id + id で効率よく1件取得します。
func (r *ChatScyllaRepository) GetChat(ctx context.Context, roomID uuid.UUID, chatID uuid.UUID) (*Chat, error) {
	var (
		gID, gRoomID, gSenderID       gocql.UUID
		message                       string
		version, createdAt, updatedAt int64
	)
	err := r.session.Query(`
		SELECT id, room_id, sender_id, message, version, created_at, updated_at
		FROM chat.chat_messages
		WHERE room_id = ? AND id = ?`,
		toGocql(roomID), toGocql(chatID),
	).WithContext(ctx).Scan(&gID, &gRoomID, &gSenderID, &message, &version, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	return &Chat{
		ID:        fromGocql(gID),
		RoomID:    fromGocql(gRoomID),
		SenderID:  fromGocql(gSenderID),
		Message:   message,
		Version:   version,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (r *ChatScyllaRepository) GetChatsByRoomID(ctx context.Context, roomID uuid.UUID, limit int) ([]*Chat, error) {
	iter := r.session.Query(`
		SELECT id, room_id, sender_id, message, version, created_at, updated_at
		FROM chat.chat_messages
		WHERE room_id = ?
		LIMIT ?`,
		toGocql(roomID), limit,
	).WithContext(ctx).Iter()

	var result []*Chat
	for {
		var (
			gID, gRoomID, gSenderID       gocql.UUID
			msg                           string
			version, createdAt, updatedAt int64
		)
		if !iter.Scan(&gID, &gRoomID, &gSenderID, &msg, &version, &createdAt, &updatedAt) {
			break
		}
		result = append(result, &Chat{
			ID:        fromGocql(gID),
			RoomID:    fromGocql(gRoomID),
			SenderID:  fromGocql(gSenderID),
			Message:   msg,
			Version:   version,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ChatScyllaRepository) UpdateChat(ctx context.Context, m *Chat) error {
	return r.session.Query(`
		UPDATE chat.chat_messages
		SET message = ?, version = ?, updated_at = ?
		WHERE room_id = ? AND id = ?`,
		m.Message, m.Version, m.UpdatedAt,
		toGocql(m.RoomID), toGocql(m.ID),
	).WithContext(ctx).Exec()
}

func (r *ChatScyllaRepository) DeleteChat(ctx context.Context, roomID uuid.UUID, chatID uuid.UUID) error {
	return r.session.Query(`
		DELETE FROM chat.chat_messages
		WHERE room_id = ? AND id = ?`,
		toGocql(roomID), toGocql(chatID),
	).WithContext(ctx).Exec()
}

// history

func (r *ChatScyllaRepository) GetHistory(ctx context.Context, roomID uuid.UUID, chatID uuid.UUID) ([]*Chat, error) {
	// ToDo status が deleted 無者は除外
	iter := r.session.Query(`
		SELECT id, room_id, sender_id, message, version, created_at, updated_at
		FROM chat.chat_message_histories
		WHERE room_id = ? AND id = ?`,
		toGocql(roomID), toGocql(chatID),
	).WithContext(ctx).Iter()

	var result []*Chat
	for {
		var (
			gID, gRoomID, gSenderID       gocql.UUID
			msg                           string
			version, createdAt, updatedAt int64
		)
		if !iter.Scan(&gID, &gRoomID, &gSenderID, &msg, &version, &createdAt, &updatedAt) {
			break
		}
		result = append(result, &Chat{
			ID:        fromGocql(gID),
			RoomID:    fromGocql(gRoomID),
			SenderID:  fromGocql(gSenderID),
			Message:   msg,
			Version:   version,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ChatScyllaRepository) InsertHistory(ctx context.Context, m *Chat, status Status) error {
	return r.session.Query(`
		INSERT INTO chat.chat_message_histories
			(id, sender_id, message, version, room_id, created_at, updated_at, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		toGocql(m.ID), toGocql(m.SenderID), m.Message, m.Version,
		toGocql(m.RoomID), m.CreatedAt, m.UpdatedAt, status,
	).WithContext(ctx).Exec()
}

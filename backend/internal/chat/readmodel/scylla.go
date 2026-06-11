package readmodel

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type ChatScyllaRepository struct {
	session *gocql.Session
}

func NewChatScyllaRepository(hosts ...string) (*ChatScyllaRepository, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = "chat"
	cluster.Consistency = gocql.LocalQuorum

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create scylla session: %w", err)
	}
	return &ChatScyllaRepository{session: session}, nil
}

func (r *ChatScyllaRepository) Close() {
	r.session.Close()
}

func (r *ChatScyllaRepository) InsertChat(ctx context.Context, m *ChatReadModel) error {
	batch := r.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	batch.Query(`
		INSERT INTO chat.chat_messages
			(room_id, created_at, id, sender_id, sender_name, sender_icon, message, version, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.RoomID, m.CreatedAt, m.ID,
		m.Sender.ID, m.Sender.Name, m.Sender.Icon,
		m.Message, m.Version, m.UpdatedAt,
	)
	batch.Query(`
		INSERT INTO chat.chat_messages_by_id (id, room_id, created_at)
		VALUES (?, ?, ?)`,
		m.ID, m.RoomID, m.CreatedAt,
	)
	return r.session.ExecuteBatch(batch)
}

// lookupByID は id → (room_id, created_at) をルックアップテーブルから取得します。
func (r *ChatScyllaRepository) lookupByID(ctx context.Context, chatID uuid.UUID) (roomID uuid.UUID, createdAt int64, err error) {
	err = r.session.Query(`
		SELECT room_id, created_at FROM chat.chat_messages_by_id WHERE id = ?`,
		chatID,
	).WithContext(ctx).Scan(&roomID, &createdAt)
	return
}

// GetChatByID は id でルックアップテーブルを引いてから主テーブルを取得します（ALLOW FILTERING 不使用）。
func (r *ChatScyllaRepository) GetChatByID(ctx context.Context, chatID uuid.UUID) (*ChatReadModel, error) {
	roomID, createdAt, err := r.lookupByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	var (
		id, senderID           uuid.UUID
		senderName, senderIcon, message string
		version, updatedAt              int64
	)
	err = r.session.Query(`
		SELECT id, room_id, sender_id, sender_name, sender_icon, message, version, created_at, updated_at
		FROM chat.chat_messages
		WHERE room_id = ? AND created_at = ? AND id = ?`,
		roomID, createdAt, chatID,
	).WithContext(ctx).Scan(
		&id, &roomID, &senderID, &senderName, &senderIcon,
		&message, &version, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ChatReadModel{
		ID:        id,
		RoomID:    roomID,
		Sender:    User{ID: senderID, Name: senderName, Icon: senderIcon},
		Message:   message,
		Version:   version,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (r *ChatScyllaRepository) UpdateChat(ctx context.Context, m *ChatReadModel) error {
	return r.session.Query(`
		UPDATE chat.chat_messages
		SET message = ?, version = ?, updated_at = ?
		WHERE room_id = ? AND created_at = ? AND id = ?`,
		m.Message, m.Version, m.UpdatedAt,
		m.RoomID, m.CreatedAt, m.ID,
	).WithContext(ctx).Exec()
}

// DeleteChat は id を起点にルックアップし、両テーブルを BATCH DELETE します。
func (r *ChatScyllaRepository) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	roomID, createdAt, err := r.lookupByID(ctx, chatID)
	if err != nil {
		return err
	}

	batch := r.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	batch.Query(`
		DELETE FROM chat.chat_messages
		WHERE room_id = ? AND created_at = ? AND id = ?`,
		roomID, createdAt, chatID,
	)
	batch.Query(`
		DELETE FROM chat.chat_messages_by_id WHERE id = ?`,
		chatID,
	)
	return r.session.ExecuteBatch(batch)
}

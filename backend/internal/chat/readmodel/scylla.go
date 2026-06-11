package readmodel

import (
	"context"
	"fmt"
	"log"
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

// toGocql は uuid.UUID を gocql.UUID に変換します（両者とも [16]byte）。
func toGocql(u uuid.UUID) gocql.UUID { return gocql.UUID(u) }

// fromGocql は gocql.UUID を uuid.UUID に変換します。
func fromGocql(u gocql.UUID) uuid.UUID { return uuid.UUID(u) }

func (r *ChatScyllaRepository) InsertChat(ctx context.Context, m *ChatReadModel) error {
	batch := r.session.NewBatch(gocql.LoggedBatch).WithContext(ctx)
	batch.Query(`
		INSERT INTO chat.chat_messages
			(room_id, created_at, id, sender_id, sender_name, sender_icon, message, version, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		toGocql(m.RoomID), m.CreatedAt, toGocql(m.ID),
		toGocql(m.Sender.ID), m.Sender.Name, m.Sender.Icon,
		m.Message, m.Version, m.UpdatedAt,
	)
	batch.Query(`
		INSERT INTO chat.chat_messages_by_id (id, room_id, created_at)
		VALUES (?, ?, ?)`,
		toGocql(m.ID), toGocql(m.RoomID), m.CreatedAt,
	)
	log.Println("Inserted chat.chat_messages_by_id:", toGocql(m.ID))
	return r.session.ExecuteBatch(batch)
}

// lookupByID は id → (room_id, created_at) をルックアップテーブルから取得します。
func (r *ChatScyllaRepository) lookupByID(ctx context.Context, chatID uuid.UUID) (roomID uuid.UUID, createdAt int64, err error) {
	var gRoomID gocql.UUID
	err = r.session.Query(`
		SELECT room_id, created_at FROM chat.chat_messages_by_id WHERE id = ?`,
		toGocql(chatID),
	).WithContext(ctx).Scan(&gRoomID, &createdAt)
	roomID = fromGocql(gRoomID)
	return
}

// GetChatByID は id でルックアップテーブルを引いてから主テーブルを取得します（ALLOW FILTERING 不使用）。
func (r *ChatScyllaRepository) GetChatByID(ctx context.Context, chatID uuid.UUID) (*ChatReadModel, error) {
	roomID, createdAt, err := r.lookupByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	var (
		gID, gRoomID, gSenderID         gocql.UUID
		senderName, senderIcon, message string
		version, updatedAt              int64
	)
	err = r.session.Query(`
		SELECT id, room_id, sender_id, sender_name, sender_icon, message, version, created_at, updated_at
		FROM chat.chat_messages
		WHERE room_id = ? AND created_at = ? AND id = ?`,
		toGocql(roomID), createdAt, toGocql(chatID),
	).WithContext(ctx).Scan(
		&gID, &gRoomID, &gSenderID, &senderName, &senderIcon,
		&message, &version, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ChatReadModel{
		ID:        fromGocql(gID),
		RoomID:    fromGocql(gRoomID),
		Sender:    User{ID: fromGocql(gSenderID), Name: senderName, Icon: senderIcon},
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
		toGocql(m.RoomID), m.CreatedAt, toGocql(m.ID),
	).WithContext(ctx).Exec()
}

func (r *ChatScyllaRepository) GetChatsByRoomID(ctx context.Context, roomID uuid.UUID, limit int) ([]*ChatReadModel, error) {
	iter := r.session.Query(`
		SELECT id, room_id, sender_id, sender_name, sender_icon, message, version, created_at, updated_at
		FROM chat.chat_messages
		WHERE room_id = ?
		LIMIT ?`,
		toGocql(roomID), limit,
	).WithContext(ctx).Iter()

	var result []*ChatReadModel
	for {
		var (
			gID, gRoomID, gSenderID       gocql.UUID
			senderName, senderIcon, msg   string
			version, createdAt, updatedAt int64
		)
		if !iter.Scan(&gID, &gRoomID, &gSenderID, &senderName, &senderIcon, &msg, &version, &createdAt, &updatedAt) {
			break
		}
		log.Println("ScyllaDB returned chat:", fromGocql(gID), "room:", fromGocql(gRoomID), "sender:", fromGocql(gSenderID))
		result = append(result, &ChatReadModel{
			ID:        fromGocql(gID),
			RoomID:    fromGocql(gRoomID),
			Sender:    User{ID: fromGocql(gSenderID), Name: senderName, Icon: senderIcon},
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
		toGocql(roomID), createdAt, toGocql(chatID),
	)
	batch.Query(`
		DELETE FROM chat.chat_messages_by_id WHERE id = ?`,
		toGocql(chatID),
	)
	return r.session.ExecuteBatch(batch)
}

package domain

import (
	"github.com/google/uuid"
)

type Chat struct {
	ID       uuid.UUID
	SenderID uuid.UUID
	RoomID   uuid.UUID
	Message  string
	Version  int64
}

func NewChat(senderId uuid.UUID, roomId uuid.UUID, message string) *Chat {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return &Chat{
		ID:       id,
		SenderID: senderId,
		RoomID:   roomId,
		Message:  message,
		Version:  1,
	}
}

func NewEditedChat(id uuid.UUID, senderId uuid.UUID, roomId uuid.UUID, message string) *Chat {
	return &Chat{
		ID:       id,
		SenderID: senderId,
		RoomID:   roomId,
		Message:  message,
		Version:  1, // ToDo
	}
}

//func (c *Chat) AsEditEvent() ChatEvent {
//	jsonChat, err := json.Marshal(c)
//	if err != nil {
//		panic(err)
//	}
//	return ChatEvent{
//		ChatID:    c.ID,
//		EventType: ChatEditedPayload,
//		Version:   c.Version,
//		Payload:   jsonChat,
//		Timestamp: time.Now().Unix(),
//	}
//}
//
//func (c *Chat) AsDeleteEvent() ChatEvent {
//	jsonChat, err := json.Marshal(c)
//	if err != nil {
//		panic(err)
//	}
//	return ChatEvent{
//		ChatID:    c.ID,
//		EventType: ChatDeletedEvent,
//		Version:   c.Version,
//		Payload:   jsonChat,
//		Timestamp: time.Now().Unix(),
//	}
//}
//
//func AsOutbox(eventId int64, e ChatEvent) ChatEventOutbox {
//	marshaledEvent, err := json.Marshal(e)
//	if err != nil {
//		panic(err)
//	}
//	return ChatEventOutbox{
//		EventId:   eventId,
//		EventType: e.EventType,
//		Payload:   marshaledEvent,
//		Timestamp: e.Timestamp,
//	}
//}

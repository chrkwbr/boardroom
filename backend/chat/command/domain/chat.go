package domain

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID
	Sender    User
	Room      Room
	Message   string
	Version   int64
	Timestamp int64
}

func NewChat(sender User, room Room, message string) *Chat {
	return &Chat{
		ID:        uuid.New(),
		Sender:    sender,
		Room:      room,
		Message:   message,
		Version:   1,
		Timestamp: time.Now().Unix(),
	}
}

//func (c *Chat) Edit(message string) Chat {
//	return Chat{
//		ID:        c.ID,
//		Sender:    c.Sender,
//		Room:      c.Room,
//		Message:   message,
//		Version:   c.Version + 1,
//		Timestamp: c.Timestamp,
//	}
//
//}

//type ChatEventOutbox struct {
//	EventId   int64
//	EventType string
//	Payload   []byte
//	Timestamp int64
//}

func (c *Chat) NewCreatedEvent() *ChatEvent {
	//jsonChat, err := json.Marshal(c)
	//if err != nil {
	//	panic(err)
	//}
	return &ChatEvent{
		Type:       EventTypeCreated,
		OccurredAt: c.Timestamp,
		Payload: ChatCreatedEvent{
			ID:       c.ID,
			RoomID:   c.Room.ID,
			SenderID: c.Sender.ID,
			Message:  c.Message,
			Version:  c.Version,
		},
	}
}

//func (c *Chat) AsEditEvent() ChatEvent {
//	jsonChat, err := json.Marshal(c)
//	if err != nil {
//		panic(err)
//	}
//	return ChatEvent{
//		ChatId:    c.ID,
//		EventType: ChatEditedEvent,
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
//		ChatId:    c.ID,
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

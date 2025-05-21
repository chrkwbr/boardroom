package handler

type WsChatEvent struct {
	ChatId    string `json:"id"`
	EventType string `json:"event_type"`
	Sender    string `json:"sender"`
	Room      string `json:"room"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

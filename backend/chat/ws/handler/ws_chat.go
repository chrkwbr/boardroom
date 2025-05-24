package handler

type WsChatEvent struct {
	ChatId    string `json:"id"`
	EventType string `json:"event_type"`
	Sender    string `json:"sender"`
	Room      string `json:"room"`
	Message   string `json:"message"`
	Version   int64  `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

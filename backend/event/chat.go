package event

type ChatEvent struct {
	EventType string `json:"event_type"`
	Message   string `json:"message"`
}

type ChatMessage struct {
	ID        string `json:"id"`
	Sender    string `json:"sender"`
	Room      string `json:"room"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

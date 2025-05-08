package hub

type Client struct {
	send chan []byte
}

func NewClient(bufferSIze int32) *Client {
	return &Client{
		send: make(chan []byte, bufferSIze),
	}
}

func (c *Client) Receive(f func([]byte)) {
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			f(msg)
		}
	}
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) CreateAndRegisterClient(bufferSIze int32) *Client {
	client := NewClient(bufferSIze)
	h.registerClient(client)
	return client
}

func (h *Hub) registerClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

func (h *Hub) BroadcastMessage(message []byte) {
	h.broadcast <- message
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

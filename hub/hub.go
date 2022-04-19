package hub

import (
  "github.com/google/uuid"
  "github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

func (h *Hub) NewClient(conn *websocket.Conn, ipAddress string, name string) (*Client) {
  id := uuid.New().String()
  if name == "" {
    name = id
  }
  cl := &Client{hub: h, conn: conn, send: make(chan []byte, 256), Name: name, IPAddress: ipAddress, Identifier: id}

  // Allow collection of memory referenced by the caller by doing all work in
  // new goroutines.
  go cl.writePump()
  go cl.readPump()

  // use channel so thread safe?
  h.register <- cl

  return cl
}

func (h *Hub) BrodcastClients() {
	for _, client := range h.clients {
    client.SendClientList()
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.Identifier] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.Identifier]; ok {
				delete(h.clients, client.Identifier)
				close(client.send)
			}
		}
		h.BrodcastClients()
	}
}

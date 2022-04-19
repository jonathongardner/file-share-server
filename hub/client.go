package hub

import (
	"log"
	"time"

  "github.com/gorilla/websocket"
)


const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2048
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Name of this user (can be empty)
	// IPAddress of client
	// Identifier: Unique
	Name       string `json:"name"`
	IPAddress  string `json:"ipAddress"`
	Identifier string `json:"identifier"`
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
    c.Delete()
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, messageByte, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v\n", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		m := unmarshalMessage(messageByte)
		switch m.Type {
		case "AddName":
			c.Name = m.Data.(string)
			c.hub.BrodcastClients()
		case "SDP":
			log.Println("SDP")
      data := m.Data.(map[string]interface {})
			identifier := data["identifier"]
      peer, ok := c.hub.clients[identifier.(string)]
      // if peer not found just ignore for now
      if !ok {
        return
      }
      peer.SendSDP(c, data["description"])
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) Delete() {
	c.hub.unregister <- c
}

func (c *Client) sendMessage(tpe string, data interface{}) {
  m := &Message{Type: tpe, Data: data}

  select {
  case c.send <- m.Json():
  default:
    c.Delete()
  }
}

func (c *Client) SendSDP(fromClient *Client, description interface{}) {
  data := map[string]interface{}{"client": fromClient, "description": description}
  c.sendMessage("SDP", data)
}

func (c *Client) SendClientList() {
  h := c.hub
  clients := make([]*Client, len(h.clients) - 1)

  i := 0
  for id, client := range h.clients {
    if id != c.Identifier {
      clients[i] = client
      i++
    }
  }

  c.sendMessage("ListClients", clients)
}

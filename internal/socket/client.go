package socket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgtype"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID pgtype.UUID
	game   *Game
}

// writePump pushes messages from c.send channel to the client
func (c *Client) writePump() {
	defer c.conn.Close()
	for {
		message, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
}

// readPump listens for messages from the client
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var req WSRequest
		if err := json.Unmarshal(message, &req); err != nil {
			log.Printf("Ошибка парсинга конверта: %v", err)
			continue
		}

		c.hub.incoming <- MessageEvent{
			Sender:  c,
			Payload: req,
		}
	}
}

func (c *Client) sendJSON(msgType string, payload WSPayload) {
	response := WSResponse{
		Type:    msgType,
		Payload: payload,
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Ошибка маршалинга JSON: %v", err)
		return
	}

	c.send <- data
}

package socket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minakdanCVUT/GoChess/internal/security"
	"github.com/minakdanCVUT/GoChess/internal/service"
)

type MessageEvent struct {
	Sender  *Client
	Payload WSRequest
}

type Hub struct {
	games         map[pgtype.UUID]*Game
	waitingPlayer *Client
	register      chan *Client
	unregister    chan *Client
	incoming      chan MessageEvent
	service       *service.GameService
}

func NewHub(s *service.GameService) *Hub {
	return &Hub{
		games:      make(map[pgtype.UUID]*Game),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		incoming:   make(chan MessageEvent, 256),
		service:    s,
	}
}

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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// extract userID from context, that AuthMiddleware put in manually from url query
	userID, err := security.ExtractUserIDFromContext(r.Context())
	if err != nil {
		return
	}
	// upgrade http request to websocket, that came like http handshake
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// create a new client for websocket connection with an empty game
	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
		game:   nil,
	}
	// push new client to the hub register channel
	client.hub.register <- client
	// turn on two different gorutines for write pump(channel) and read pump
	go client.writePump()
	go client.readPump()
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

func UnmarshalAsType[T WSPayload](raw json.RawMessage) (T, error) {
	var zero T
	var data T
	if err := json.Unmarshal(raw, &data); err != nil {
		return zero, err
	}
	return data, nil
}

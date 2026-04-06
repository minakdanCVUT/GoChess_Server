package socket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/minakdanCVUT/GoChess/internal/security"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// extract userID from context, that AuthMiddleware put in manually from url query
	userID, err := security.ExtractUserIDFromContext(r.Context())
	if err != nil {
		return
	}
	// extract gameType from context, that QueryMiddleware put in manually from url query
	gameType, err := security.ExtractGameTypeFromContext(r.Context())
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

	request := GameTypeRequest{
		Sender:   client,
		GameType: security.GetGameTypeFromString(gameType),
	}

	go client.writePump()
	go client.readPump()

	// turn on two different gorutines for write pump(channel) and read pump
	client.hub.register <- request
}

func UnmarshalAsType[T WSPayload](raw json.RawMessage) (T, error) {
	var zero T
	var data T
	if err := json.Unmarshal(raw, &data); err != nil {
		return zero, err
	}
	return data, nil
}

package socket

import (
	"log"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minakdanCVUT/GoChess/internal/security"
	"github.com/minakdanCVUT/GoChess/internal/service"
)

type Hub struct {
	workers    map[security.GameType]*MatchMaker
	userRoutes map[pgtype.UUID]*MatchMaker

	register   chan GameTypeRequest
	unregister chan *Client

	incoming chan MessageEvent
	done     chan *Client
}

func NewHub(s *service.GameService) *Hub {
	h := &Hub{
		workers:    make(map[security.GameType]*MatchMaker),
		userRoutes: make(map[pgtype.UUID]*MatchMaker),
		register:   make(chan GameTypeRequest),
		unregister: make(chan *Client),
		done:       make(chan *Client),
		incoming:   make(chan MessageEvent, 256),
	}

	for _, gt := range security.AllGameTypes() {
		worker := NewMatchMaker(h, s, gt)
		h.workers[gt] = worker
		go worker.Run()
	}

	return h
}

func (h *Hub) Run() {
	for {
		select {
		case req := <-h.register:
			userID := req.Sender.userID

			if _, alreadyBusy := h.userRoutes[userID]; alreadyBusy {
				//req.Sender.sendJSON("error", "Вы уже находитесь в очереди или в игре")
				continue
			}

			if worker, ok := h.workers[req.GameType]; ok {
				h.userRoutes[userID] = worker
				worker.Register <- req.Sender
			}

		case client := <-h.unregister:
			if worker, ok := h.userRoutes[client.userID]; ok {
				worker.Unregister <- client
				delete(h.userRoutes, client.userID)
			}

		case event := <-h.incoming:
			if worker, ok := h.userRoutes[event.Sender.userID]; ok {
				worker.Input <- event
			}
		case client := <-h.done:
			log.Printf("[hub] player %s is now free", client.userID.String())
			delete(h.userRoutes, client.userID)
		}

	}
}

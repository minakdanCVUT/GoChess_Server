package socket

import (
	"log"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minakdanCVUT/GoChess/internal/security"
	"github.com/minakdanCVUT/GoChess/internal/service"
)

type MatchMaker struct {
	hub           *Hub
	games         map[pgtype.UUID]*Game
	waitingPlayer *Client
	mode          security.GameType

	Register   chan *Client
	Unregister chan *Client

	Input   chan MessageEvent
	service *service.GameService
}

func NewMatchMaker(h *Hub, s *service.GameService, m security.GameType) *MatchMaker {
	return &MatchMaker{
		hub:     h,
		games:   make(map[pgtype.UUID]*Game),
		mode:    m,
		service: s,

		Register:   make(chan *Client),
		Unregister: make(chan *Client),

		Input: make(chan MessageEvent, 256),
	}
}

func (mm *MatchMaker) Run() {
	for {
		select {
		case client := <-mm.Register:
			if mm.waitingPlayer == nil {
				mm.waitingPlayer = client
				log.Printf("[hub] player %s connected → waiting queue", client.userID.String())
				client.sendJSON("in_queue", InQueuePayload{Message: "You are in queue. Waiting for opponent..."})
			} else {
				player1 := mm.waitingPlayer
				player2 := client
				log.Printf("[hub] player %s connected → opponent found (%s), creating game...", player2.userID.String(), player1.userID.String())
				createdGame, err := mm.service.CreateGame(player1.userID, player2.userID)
				if err != nil {
					log.Printf("[hub] failed to create game: %v", err)
					continue
				}
				newGame := &Game{
					ID:          createdGame.ID,
					WhitePlayer: player1,
					BlackPlayer: player2,

					WhiteID: player1.userID,
					BlackID: player2.userID,
					Turn:    "white",
				}
				player1.game = newGame
				player2.game = newGame

				mm.games[newGame.ID] = newGame
				log.Printf("[hub] game started | id=%s | white=%s | black=%s", newGame.ID.String(), player1.userID.String(), player2.userID.String())

				player1.sendJSON("game_started", GameStartedPayload{
					GameID: newGame.ID.String(),
					Color:  "white",
				})

				player2.sendJSON("game_started", GameStartedPayload{
					GameID: newGame.ID.String(),
					Color:  "black",
				})
				mm.waitingPlayer = nil
			}

		case client := <-mm.Unregister:
			mm.HandleDisconnect(client)

		case event := <-mm.Input:
			switch event.Payload.Type {
			case "move":
				data, err := UnmarshalAsType[MovePayload](event.Payload.Payload)
				if err != nil {
					continue
				}
				mm.HandleMove(event.Sender, &data)
			case "leave_game":
				_, err := UnmarshalAsType[LeaveGamePayload](event.Payload.Payload)
				if err != nil {
					continue
				}
				mm.HandleLeaveGame(event.Sender)
			}
		}
	}
}

func (mm *MatchMaker) HandleMove(sender *Client, data *MovePayload) {
	game := sender.game
	if game == nil {
		return
	}

	isWhiteSender := (sender == game.WhitePlayer)
	if (game.Turn == "white" && !isWhiteSender) || (game.Turn == "black" && isWhiteSender) {
		return
	}

	opponent := GetOpponent(sender)
	if opponent != nil {
		if game.Turn == "white" {
			game.Turn = "black"
		} else {
			game.Turn = "white"
		}
		opponent.sendJSON("move", data)
	}
}

func (mm *MatchMaker) HandleLeaveGame(sender *Client) {
	game := sender.game
	if game == nil {
		return
	}

	opponent := GetOpponent(sender)
	if opponent != nil {
		opponent.sendJSON("win_leave", WinCauseLeavePayload{
			Message: "You win! Your opponent left the game.",
		})
	}

	mm.closeGame(game)
}

func (mm *MatchMaker) HandleDisconnect(client *Client) {
	log.Printf("[hub] player %s disconnected", client.userID.String())
	if mm.waitingPlayer == client {
		mm.waitingPlayer = nil
		log.Printf("[hub] player %s removed from queue", client.userID.String())
		return
	}

	game := client.game
	if game != nil {
		opponent := GetOpponent(client)
		if opponent != nil {
			opponent.sendJSON("win_disconnect", WinCauseDisconnectPayload{
				Message: "You win! Opponent disconnected.",
			})
		}
		mm.closeGame(game)
	}
}

func (mm *MatchMaker) closeGame(game *Game) {
	if game.WhitePlayer != nil {
		mm.hub.done <- game.WhitePlayer
		game.WhitePlayer.game = nil
	}
	if game.BlackPlayer != nil {
		mm.hub.done <- game.BlackPlayer
		game.BlackPlayer.game = nil
	}
	delete(mm.games, game.ID)
}

func GetOpponent(sender *Client) *Client {
	game := sender.game
	if game == nil {
		return nil
	}

	if sender == game.WhitePlayer {
		return game.BlackPlayer
	}
	return game.WhitePlayer
}

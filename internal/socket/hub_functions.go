package socket

import "log"

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if h.waitingPlayer == nil {
				h.waitingPlayer = client
				log.Printf("[hub] player %s connected → waiting queue", client.userID.String())
				client.sendJSON("in_queue", InQueuePayload{Message: "You are in queue. Waiting for opponent..."})
			} else {
				player1 := h.waitingPlayer
				player2 := client
				log.Printf("[hub] player %s connected → opponent found (%s), creating game...", player2.userID.String(), player1.userID.String())
				createdGame, err := h.service.CreateGame(player1.userID, player2.userID)
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

				h.games[newGame.ID] = newGame
				log.Printf("[hub] game started | id=%s | white=%s | black=%s", newGame.ID.String(), player1.userID.String(), player2.userID.String())

				player1.sendJSON("game_started", GameStartedPayload{
					GameID: newGame.ID.String(),
					Color:  "white",
				})

				player2.sendJSON("game_started", GameStartedPayload{
					GameID: newGame.ID.String(),
					Color:  "black",
				})
				h.waitingPlayer = nil
			}
		case client := <-h.unregister:
			h.HandleDisconnect(client)
		case event := <-h.incoming:
			switch event.Payload.Type {
			case "move":
				data, err := UnmarshalAsType[MovePayload](event.Payload.Payload)
				if err != nil {
					continue
				}
				h.HandleMove(event.Sender, &data)
			case "leave_game":
				_, err := UnmarshalAsType[LeaveGamePayload](event.Payload.Payload)
				if err != nil {
					continue
				}
				h.HandleLeaveGame(event.Sender)
			}
		}
	}
}

func (h *Hub) HandleMove(sender *Client, data *MovePayload) {
	game := sender.game
	if game == nil {
		return
	}

	isWhiteSender := (sender == game.WhitePlayer)
	if (game.Turn == "white" && !isWhiteSender) || (game.Turn == "black" && isWhiteSender) {
		return
	}

	opponent := h.getOpponent(sender)
	if opponent != nil {
		if game.Turn == "white" {
			game.Turn = "black"
		} else {
			game.Turn = "white"
		}
		opponent.sendJSON("move", data)
	}
}

func (h *Hub) HandleLeaveGame(sender *Client) {
	game := sender.game
	if game == nil {
		return
	}

	opponent := h.getOpponent(sender)
	if opponent != nil {
		opponent.sendJSON("win_leave", WinCauseLeavePayload{
			Message: "You win! Your opponent left the game.",
		})
	}

	h.closeGame(game)
}

func (h *Hub) HandleDisconnect(client *Client) {
	log.Printf("[hub] player %s disconnected", client.userID.String())
	if h.waitingPlayer == client {
		h.waitingPlayer = nil
		log.Printf("[hub] player %s removed from queue", client.userID.String())
		return
	}

	game := client.game
	if game != nil {
		opponent := h.getOpponent(client)
		if opponent != nil {
			opponent.sendJSON("win_disconnect", WinCauseDisconnectPayload{
				Message: "You win! Opponent disconnected.",
			})
		}
		h.closeGame(game)
	}
}

func (h *Hub) closeGame(game *Game) {
	if game.WhitePlayer != nil {
		game.WhitePlayer.game = nil
	}
	if game.BlackPlayer != nil {
		game.BlackPlayer.game = nil
	}
	delete(h.games, game.ID)
}

func (h *Hub) getOpponent(sender *Client) *Client {
	game := sender.game
	if game == nil {
		return nil
	}

	if sender == game.WhitePlayer {
		return game.BlackPlayer
	}
	return game.WhitePlayer
}

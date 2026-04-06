package socket

import "github.com/minakdanCVUT/GoChess/internal/security"

type GameTypeRequest struct {
	Sender   *Client
	GameType security.GameType
}

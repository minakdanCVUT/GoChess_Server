package socket

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Game struct {
	ID          pgtype.UUID
	BlackPlayer *Client
	WhitePlayer *Client

	BlackID pgtype.UUID
	WhiteID pgtype.UUID

	Turn string
}

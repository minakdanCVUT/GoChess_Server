package socket

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type MessageEvent struct {
	Sender  *Client
	Payload WSRequest
}

type Game struct {
	ID          pgtype.UUID
	BlackPlayer *Client
	WhitePlayer *Client

	BlackID pgtype.UUID
	WhiteID pgtype.UUID

	Turn string
}

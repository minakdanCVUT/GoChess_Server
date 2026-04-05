package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minakdanCVUT/GoChess/internal/db"
)

type GameService struct {
	queries *db.Queries
}

func NewGameService(q *db.Queries) *GameService {
	return &GameService{queries: q}
}

func (s *GameService) CreateGame(blackId pgtype.UUID, whiteId pgtype.UUID) (*db.Game, error) {
	var params db.CreateGameParams
	params.PlayerBlackID = blackId
	params.PlayerWhiteID = whiteId
	game, err := s.queries.CreateGame(context.Background(), params)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

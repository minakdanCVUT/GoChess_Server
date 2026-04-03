-- +goose Up
CREATE TYPE IF NOT EXISTS game_status AS ENUM ('waiting', 'ongoing', 'finished', 'abandoned');
CREATE TYPE IF NOT EXISTS end_reason  AS ENUM ('checkmate', 'resign', 'timeout', 'draw_agreement', 'stalemate');

CREATE TABLE IF NOT EXISTS games (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    player_white_id UUID        NOT NULL REFERENCES users(id),
    player_black_id UUID        NOT NULL REFERENCES users(id),
    status          game_status NOT NULL DEFAULT 'waiting',
    winner_id       UUID        REFERENCES users(id),
    end_reason      end_reason,
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_games_player_white ON games(player_white_id);
CREATE INDEX IF NOT EXISTS idx_games_player_black ON games(player_black_id);

-- +goose Down
DROP INDEX IF EXISTS idx_games_player_black;
DROP INDEX IF EXISTS idx_games_player_white;
DROP TABLE IF EXISTS games;
DROP TYPE IF EXISTS end_reason;
DROP TYPE IF EXISTS game_status;

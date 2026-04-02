CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE game_status AS ENUM ('waiting', 'ongoing', 'finished', 'abandoned');
CREATE TYPE end_reason AS ENUM ('checkmate', 'resign', 'timeout', 'draw_agreement', 'stalemate');

CREATE TABLE users (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name      TEXT        NOT NULL,
    last_name       TEXT        NOT NULL,
    username        TEXT        NOT NULL UNIQUE,
    email           TEXT        NOT NULL UNIQUE,
    email_verified  BOOLEAN     NOT NULL DEFAULT FALSE,
    email_verified_at TIMESTAMPTZ,
    password TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE games (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    player_white_id  UUID        NOT NULL REFERENCES users(id),
    player_black_id  UUID        NOT NULL REFERENCES users(id),
    status           game_status NOT NULL DEFAULT 'waiting',
    winner_id        UUID        REFERENCES users(id),
    end_reason       end_reason,
    started_at       TIMESTAMPTZ,
    finished_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE moves (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id     UUID        NOT NULL REFERENCES games(id),
    player_id   UUID        NOT NULL REFERENCES users(id),
    move_number INT         NOT NULL,
    from_x      SMALLINT    NOT NULL,
    from_y      SMALLINT    NOT NULL,
    to_x        SMALLINT    NOT NULL,
    to_y        SMALLINT    NOT NULL,
    piece       TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (game_id, move_number)
);

CREATE INDEX idx_games_player_white ON games(player_white_id);
CREATE INDEX idx_games_player_black ON games(player_black_id);
CREATE INDEX idx_moves_game_id      ON moves(game_id);
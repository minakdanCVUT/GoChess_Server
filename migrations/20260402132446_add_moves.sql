-- +goose Up
CREATE TABLE IF NOT EXISTS moves (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id     UUID        NOT NULL REFERENCES games(id) ON DELETE CASCADE,
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

CREATE INDEX IF NOT EXISTS idx_moves_game_id ON moves(game_id);

-- +goose Down
DROP INDEX IF EXISTS idx_moves_game_id;
DROP TABLE IF EXISTS moves;

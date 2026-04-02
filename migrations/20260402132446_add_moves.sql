-- +goose Up
CREATE TABLE moves (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id       UUID      NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    player_id     UUID      NOT NULL REFERENCES users(id),
    move_notation TEXT        NOT NULL,
    fen_after     TEXT        NOT NULL,
    move_number   INT         NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE moves;

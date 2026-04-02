-- +goose Up
CREATE TABLE games (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    white_player_id UUID      NOT NULL REFERENCES users(id),
    black_player_id UUID      REFERENCES users(id),
    status          TEXT        NOT NULL DEFAULT 'waiting',
    winner_id       UUID      REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE games;

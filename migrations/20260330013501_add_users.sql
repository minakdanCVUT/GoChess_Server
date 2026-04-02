-- +goose Up
CREATE TABLE users (
    id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name        TEXT        NOT NULL,
    last_name         TEXT        NOT NULL,
    username          TEXT        NOT NULL UNIQUE,
    email             TEXT        NOT NULL UNIQUE,
    email_verified    BOOLEAN     NOT NULL DEFAULT FALSE,
    email_verified_at TIMESTAMPTZ,
    password          TEXT        NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
-- При откате просто удаляем таблицу users
DROP TABLE users;
CREATE TABLE IF NOT EXISTS users
(
    id        SERIAL PRIMARY KEY,
    email     VARCHAR NOT NULL UNIQUE,
    pass_hash bytea NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);

CREATE TABLE IF NOT EXISTS apps
(
    id     SERIAL PRIMARY KEY,
    name   VARCHAR NOT NULL UNIQUE,
    secret VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users (
    id        UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    username  VARCHAR(15) UNIQUE NOT NULL,
    email     TEXT        UNIQUE NOT NULL,
    password  TEXT,
    google_id TEXT        UNIQUE,
    github_id TEXT        UNIQUE,
    role      TEXT        NOT NULL DEFAULT 'user'
);

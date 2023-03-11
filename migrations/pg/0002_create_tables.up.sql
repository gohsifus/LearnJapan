CREATE TABLE IF NOT EXISTS production.users
(
    id       serial PRIMARY KEY,
    login    varchar NOT NULL UNIQUE,
    password varchar NOT NULL,
    email    varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS production.cards
(
    id        serial PRIMARY KEY,
    in_japan   varchar NOT NULL,
    in_russian varchar NOT NULL,
    mark      int     NOT NULL,
    date_add   date DEFAULT NOW(),
    user_id    int REFERENCES production.users (id)
);

CREATE TABLE IF NOT EXISTS production.sessions
(
    id varchar PRIMARY KEY,
    user_id    int       NOT NULL REFERENCES production.users (id),
    expires   timestamp NOT NULL
);
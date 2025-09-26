CREATE TABLE pairing (
    round   int         NOT NULL,
    player1 varchar(36) NOT NULL,
    player2 varchar(36) NOT NULL,
    wins1   int         NOT NULL,
    wins2   int         NOT NULL,
    draws   int         NOT NULL
);

CREATE INDEX pairings_round_players_idx ON pairing (round, player1, player2);
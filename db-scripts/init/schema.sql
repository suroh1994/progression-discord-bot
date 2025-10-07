CREATE TABLE player_card_pool
(
    player_id        varchar(36)  NOT NULL,
    name             varchar(255) NOT NULL,
    set_code         varchar(4)   NOT NULL,
    collector_number int          NOT NULL,
    count            int          NOT NULL,
    PRIMARY KEY (player_id, set_code, collector_number)
);

CREATE TABLE league
(
    round      int         NOT NULL,
    active     bool        NOT NULL,
    started_at timestamptz NULL
);

CREATE TABLE pairing
(
    round      int         NOT NULL,
    player_id1 varchar(36) NOT NULL,
    player_id2 varchar(36) NOT NULL,
    wins1      int         NOT NULL,
    wins2      int         NOT NULL,
    draws      int         NOT NULL
);

CREATE INDEX pairings_round_players_idx ON pairing (round, player_id1, player_id2);

CREATE TABLE player
(
    id              varchar(36) NOT NULL,
    wild_card_count int         NOT NULL,
    wild_pack_count int         NOT NULL,
    PRIMARY KEY (id)
);
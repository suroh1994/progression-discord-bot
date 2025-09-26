CREATE TABLE player_card_pool (
    id                  varchar(36) NOT NULL,
    set_code            varchar(4)  NOT NULL,
    collector_number    int         NOT NULL,
    count               int         NOT NULL,
    PRIMARY KEY (id, set_code, collector_number)
);
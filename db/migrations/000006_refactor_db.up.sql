DROP TABLE IF EXISTS buildinglocations;
DROP TABLE IF EXISTS buildings;

ALTER TABLE IF EXISTS accountcharacters
    RENAME TO account_characters;
ALTER TABLE IF EXISTS chatmessages
    RENAME TO chat_messages;

CREATE TABLE IF NOT EXISTS towns
(
    id         serial      PRIMARY KEY,
    x          int         NOT NULL,
    y          int         NOT NULL,
    name       varchar(40) NOT NULL,
    owner_name varchar(25) NOT NULL,
    population int         NOT NULL DEFAULT 0,

    UNIQUE (x, y)
);

DROP TABLE IF EXISTS chunks;
CREATE TABLE chunks
(
    x           int NOT NULL,
    y           int NOT NULL,
    data        bytea        DEFAULT NULL,
    trees int NOT NULL DEFAULT 0,
    stones int NOT NULL DEFAULT 0,
    animals int NOT NULL DEFAULT 0,
    plants int NOT NULL DEFAULT 0,

    UNIQUE (x, y)
);


CREATE TABLE IF NOT EXISTS characters
(
    id   serial PRIMARY KEY,
    name varchar(25) NOT NULL,
    gold bigint      NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS accounts
(
    id       serial PRIMARY KEY,
    login    varchar(25) UNIQUE NOT NULL,
    password varchar(32)        NOT NULL,
    salt     varchar(10)        NOT NULL
);

CREATE TABLE IF NOT EXISTS accountCharacters
(
    account_id   int NOT NULL,
    character_id int NOT NULL
);

CREATE TABLE IF NOT EXISTS buildings
(
    id   serial PRIMARY KEY,
    name varchar(25) NOT NULL,
    cost int         NOT NULL
);

CREATE TABLE IF NOT EXISTS buildingLocations
(
    building_id int NOT NULL,
    owner_id    int NOT NULL,
    location    real ARRAY[3],
    UNIQUE (location)
);

CREATE TABLE IF NOT EXISTS chatMessages
(
    message_id  serial PRIMARY KEY,
    sender_name varchar(25)  NOT NULL,
    text        varchar(200) NOT NULL
);


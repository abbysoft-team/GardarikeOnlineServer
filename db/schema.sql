CREATE TABLE characters (
    id serial PRIMARY KEY,
    name varchar(25) NOT NULL,
    gold bigint NOT NULL DEFAULT 0
)

CREATE TABLE accounts (
    id serial PRIMARY KEY,
    login varchar(25) NOT NULL,
    password varchar(25) NOT NULL,
    salt varchar(10) NOT NULL
)

CREATE TABLE accountCharacters (
    account_id int NOT NULL,
    character_id int NOT NULL
)

ALTER TABLE characters
DROP COLUMN IF EXISTS gold;

CREATE TABLE IF NOT EXISTS resources (
    character_id int PRIMARY KEY,
    wood int NOT NULL DEFAULT 30,
    stone int NOT NULL DEFAULT 30,
    food int NOT NULL DEFAULT 50,
    leather int NOT NULL DEFAULT 10
)

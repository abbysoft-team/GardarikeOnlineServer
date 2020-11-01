ALTER TABLE characters
ADD COLUMN max_population int NOT NULL DEFAULT 0,
ADD COLUMN current_population int NOT NULL DEFAULT 0;

ALTER TABLE buildings
ADD COLUMN population_bonus int NOT NULL DEFAULT 5;

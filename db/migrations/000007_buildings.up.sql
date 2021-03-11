CREATE TABLE IF NOT EXISTS town_buildings (
   town_id int NOT NULL,
   building_id int NOT NULL,
   location_x int NOT NULL,
   location_y int NOT NULL,

   UNIQUE (town_id, location_x, location_y)
);

CREATE TABLE IF NOT EXISTS production_rates
(
    character_id int PRIMARY KEY,
    wood         int NOT NULL DEFAULT 0,
    stone        int NOT NULL DEFAULT 0,
    food         int NOT NULL DEFAULT 0,
    leather      int NOT NULL DEFAULT 0
)

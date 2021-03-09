CREATE TABLE IF NOT EXISTS town_buildings (
   town_id int NOT NULL,
   building_id int NOT NULL,
   location_x int NOT NULL,
   location_y int NOT NULL,

   UNIQUE (town_id, location_x, location_y)
);

CREATE TABLE IF NOT EXISTS chunks
(
    x int NOT NULL,
    y int NOT NULL,
    data bytea DEFAULT NULL,
    trees_count int NOT NULL DEFAULT 0,

    UNIQUE (x, y)
);

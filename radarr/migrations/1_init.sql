CREATE TABLE IF NOT EXISTS movies (
    "movie" INTEGER NOT NULL,
    "release_date" DATETIME NOT NULL,
    "search_date" DATETIME NULL,
    "type" VARCHAR(10),
    PRIMARY KEY(movie)
)
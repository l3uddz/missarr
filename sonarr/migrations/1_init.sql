CREATE TABLE IF NOT EXISTS series (
    "series" INTEGER NOT NULL,
    "season" INTEGER NOT NULL,
    "air_date" DATETIME NOT NULL,
    "search_date" DATETIME NULL,
    PRIMARY KEY(series)
)
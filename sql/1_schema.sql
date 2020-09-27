CREATE TABLE IF NOT EXISTS shortcuts (
    short_code text NOT NULL CONSTRAINT shortcuts_pk PRIMARY KEY,
    long_url text NOT NULL
);

#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    \c playlist_db

    BEGIN;

    CREATE TABLE songs
    (
        id       serial      not null unique,
        name     varchar(50) not null,
        duration int         not null
    );

    INSERT INTO songs (name, duration) VALUES ('Harry Styles - Watermelon Sugar', 188);
    INSERT INTO songs (name, duration) VALUES ('Rihanna - Desperado', 206);
    INSERT INTO songs (name, duration) VALUES ('Rihanna - Pour It Up', 196);
    INSERT INTO songs (name, duration) VALUES ('QUOK - Concorde', 196);
    INSERT INTO songs (name, duration) VALUES ('Paramore - Decode', 261);
    COMMIT;
EOSQL

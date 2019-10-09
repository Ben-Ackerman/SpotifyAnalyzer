CREATE EXTENSION pgcrypto;

GRANT ALL PRIVILEGES ON DATABASE spotify_analyzer TO db1;
CREATE SCHEMA IF NOT EXISTS info;
CREATE TABLE info.tracks(
   spotify_id TEXT PRIMARY KEY,
   track TEXT,
   artist TEXT,
   lyrics TEXT,
   geniusURI TEXT
);

GRANT ALL PRIVILEGES ON SCHEMA info TO db1;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA info TO db1;
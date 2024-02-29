CREATE TABLE IF NOT EXISTS "game_redirect" (
  "id" varchar(36) NOT NULL,
  "source_id" varchar(36) NOT NULL,
  "date_added" timestamp NOT NULL,
  PRIMARY KEY("id", "source_id")
);
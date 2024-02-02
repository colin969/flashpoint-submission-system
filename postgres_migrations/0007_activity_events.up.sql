CREATE TABLE "activity_events"
(
    "id"              serial PRIMARY KEY,
    "uid"             bigint      NOT NULL,
    "created_at"      timestamp NOT NULL,
    "event_area"      text      NOT NULL,
    "event_operation" text      NOT NULL,
    "event_data"      jsonb     NOT NULL
);
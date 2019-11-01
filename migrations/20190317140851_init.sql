-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS hstore;
CREATE TABLE IF NOT EXISTS events(
  id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
  title text,
  starts timestamp,
  ends timestamp,
  description text,
  category text,
  is_recurring bool,
  frequency text,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS events;

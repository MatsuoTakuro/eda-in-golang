-- +goose Up
CREATE TABLE stores_cache (
  id         text        NOT NULL,
  name       text        NOT NULL,
  location   text        NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TRIGGER created_at_stores_trgr
  BEFORE UPDATE
  ON stores_cache
  FOR EACH ROW
EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_stores_trgr
  BEFORE UPDATE
  ON stores_cache
  FOR EACH ROW
EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE products_cache (
  id         text        NOT NULL,
  store_id   text        NOT NULL,
  name       text        NOT NULL,
  price      decimal(9, 4) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TRIGGER created_at_products_trgr
  BEFORE UPDATE
  ON products_cache
  FOR EACH ROW
EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_products_trgr
  BEFORE UPDATE
  ON products_cache
  FOR EACH ROW
EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE shopping_lists (
  id              text        NOT NULL,
  order_id        text        NOT NULL,
  stops           bytea       NOT NULL,
  assigned_bot_id text        NOT NULL,
  status          text        NOT NULL,
  created_at      timestamptz NOT NULL DEFAULT NOW(),
  updated_at      timestamptz NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE INDEX shopping_lists_order_id_idx ON shopping_lists (order_id);
CREATE INDEX shopping_lists_availability_idx ON shopping_lists (status, created_at) WHERE status = 'available';

CREATE TRIGGER created_at_shopping_lists_trgr
  BEFORE UPDATE
  ON shopping_lists
  FOR EACH ROW
EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_shopping_lists_trgr
  BEFORE UPDATE
  ON shopping_lists
  FOR EACH ROW
EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE events (
  stream_id      text        NOT NULL,
  stream_name    text        NOT NULL,
  stream_version int         NOT NULL,
  event_id       text        NOT NULL,
  event_name     text        NOT NULL,
  event_data     bytea       NOT NULL,
  occurred_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (stream_id, stream_name, stream_version)
);

CREATE TABLE snapshots (
  stream_id      text        NOT NULL,
  stream_name    text        NOT NULL,
  stream_version int         NOT NULL,
  snapshot_name  text        NOT NULL,
  snapshot_data  bytea       NOT NULL,
  updated_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (stream_id, stream_name)
);

CREATE TRIGGER updated_at_snapshots_trgr
  BEFORE UPDATE
  ON snapshots
  FOR EACH ROW
EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE inbox (
  id          text        NOT NULL,
  name        text        NOT NULL,
  subject     text        NOT NULL,
  data        bytea       NOT NULL,
  received_at timestamptz NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE outbox (
  id           text  NOT NULL,
  name         text  NOT NULL,
  subject      text  NOT NULL,
  data         bytea NOT NULL,
  published_at timestamptz,
  PRIMARY KEY (id)
);

CREATE INDEX unpublished_idx ON outbox (published_at) WHERE published_at IS NULL;

CREATE TABLE sagas (
  id           text        NOT NULL,
  name         text        NOT NULL,
  data         bytea       NOT NULL,
  step         int         NOT NULL,
  done         bool        NOT NULL,
  compensating bool        NOT NULL,
  updated_at   timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id, name)
);

CREATE TRIGGER updated_at_sagas_trgr
  BEFORE UPDATE
  ON sagas
  FOR EACH ROW
EXECUTE PROCEDURE updated_at_trigger();

-- +goose Down
DROP SCHEMA IF EXISTS depot CASCADE;

DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;

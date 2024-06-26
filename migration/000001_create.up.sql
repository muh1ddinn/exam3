CREATE TABLE users (
  id UUID PRIMARY KEY,
  mail VARCHAR UNIQUE,
  first_name VARCHAR NOT NULL,
  last_name VARCHAR,
  password VARCHAR NOT NULL,
  phone VARCHAR UNIQUE,
  sex VARCHAR NOT NULL,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at INTEGER DEFAULT 0
);

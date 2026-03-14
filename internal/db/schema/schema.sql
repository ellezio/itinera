CREATE TABLE resources (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  source TEXT NOT NULL,
  source_type TEXT NOT NULL,
  status_id INTEGER NOT NULL
);

CREATE TABLE statuses (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE tags (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE resource_tags (
  resource_id INTEGER NOT NULL,
  tag_id INTEGER NOT NULL
);

CREATE TABLE notes (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  entity_id INTEGER NOT NULL,
  entity_type TEXT NOT NULL
);

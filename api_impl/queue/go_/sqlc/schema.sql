CREATE TABLE parigot_test_queue (
  id   INTEGER PRIMARY KEY,
  name text    NOT NULL
);

CREATE TABLE  parigot_test_queue_id_to_key  (
  id_low  UNSIGNED BIG INT,  -- works out to INTEGER
  id_high  UNSIGNED BIG INT,  -- works out to INTEGER
  queue_key INTEGER,
  PRIMARY KEY (id_low, id_high),
  FOREIGN KEY(queue_key) REFERENCES parigot_test_queue(id)
) WITHOUT ROWID;


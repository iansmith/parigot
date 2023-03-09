PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS parigot_test_queue (
  id   INTEGER PRIMARY KEY,
  name text    NOT NULL
);

CREATE TABLE IF NOT EXISTS parigot_test_queue_id_to_key  (
  id_low  UNSIGNED BIG INT,  -- works out to INTEGER
  id_high  UNSIGNED BIG INT,  -- works out to INTEGER
  queue_key INTEGER,
  PRIMARY KEY (id_low, id_high),
  
  CONSTRAINT fk_queueid
    FOREIGN KEY (queue_key)
    REFERENCES parigot_test_queue (id)
    ON DELETE CASCADE

) WITHOUT ROWID;

CREATE TABLE IF NOT EXISTS parigot_test_message  (
  id_low  UNSIGNED BIG INT,  -- works out to INTEGER
  id_high  UNSIGNED BIG INT,  -- works out to INTEGER
  queue_key INTEGER,
  last_received TEXT, -- date, when we last gave it to user code
  received_count UNSIGNED BIG INT DEFAULT 0,  -- works out to INTEGER
  original_sent TEXT DEFAULT CURRENT_TIMESTAMP, -- date, original received time
  marked_done TEXT, -- date (NULL meas not marked yet)
  sender BLOB, -- serialized protobuf or NULL
  payload BLOB NOT NULL, -- must have a payload
  
  PRIMARY KEY (id_low, id_high),
  CONSTRAINT fk_queueid_to_msg_id
    FOREIGN KEY (queue_key)
    REFERENCES parigot_test_queue (id)
    ON DELETE CASCADE
) WITHOUT ROWID;

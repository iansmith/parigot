-- name: GetQueue :one
SELECT * FROM parigot_test_queue
WHERE id = ?;

-- name: CreateQueue :one
INSERT INTO parigot_test_queue (
  name
) VALUES (
  ?
)
RETURNING *;

-- name: CreateIdToKeyMapping :one
INSERT INTO parigot_test_queue_id_to_key (
  id_low, id_high, queue_key
) VALUES (
  ?, ?, ?
)
RETURNING *;


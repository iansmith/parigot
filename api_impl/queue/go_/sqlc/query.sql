-- name: GetQueue :one
SELECT * FROM parigot_test_queue
WHERE id = ?;

-- name: getKeyFromQueueId :one
SELECT * FROM parigot_test_queue_id_to_key
WHERE id_low = ? AND id_high = ?;

-- name: DeleteQueue :exec
DELETE FROM parigot_test_queue 
WHERE id = ?;

-- name: testDestroyAll :exec
DELETE FROM parigot_test_queue;

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

-- name: CreateMessage :one
INSERT INTO parigot_test_message (
  id_low, id_high, queue_key, sender, payload
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;


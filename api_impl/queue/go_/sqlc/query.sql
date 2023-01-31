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

-- name: RetrieveMessage :many
SELECT * FROM parigot_test_message 
WHERE queue_key = ? 
ORDER BY original_sent
LIMIT 10;

-- name: UpdateMessageRetrieved :exec
UPDATE parigot_test_message 
SET last_received=now(), received_count=last_received+1
WHERE queue_key=? AND id_low=? AND id_high=?;

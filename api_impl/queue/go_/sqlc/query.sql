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

-- name: Length :one
SELECT count(*) FROM parigot_test_message 
WHERE queue_key = ? AND marked_done IS NULL;

-- name: Locate :one
SELECT parigot_test_queue_id_to_key.id_high, parigot_test_queue_id_to_key.id_low 
FROM parigot_test_queue INNER JOIN parigot_test_queue_id_to_key on  parigot_test_queue.id = parigot_test_queue_id_to_key.queue_key
WHERE parigot_test_queue.name = ? ;

-- name: allMessages :many
SELECT queue_key,id_low, id_high, marked_done, original_sent, last_received FROM parigot_test_message
WHERE queue_key = ?  AND marked_done IS NULL;

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
WHERE queue_key = ? AND marked_done IS NULL
ORDER BY original_sent
LIMIT 10;

-- name: MarkDone :exec
UPDATE parigot_test_message 
SET marked_done = 'now'
WHERE queue_key = ? AND marked_done IS NULL AND id_low=? AND id_high=?;

-- name: UpdateMessageRetrieved :exec
UPDATE parigot_test_message 
SET last_received='now', received_count=received_count+1
WHERE queue_key=? AND id_low=? AND id_high=?;

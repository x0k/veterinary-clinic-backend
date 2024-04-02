-- name: GetClient :one
SELECT * FROM appointment_clients
WHERE id = ? LIMIT 1;

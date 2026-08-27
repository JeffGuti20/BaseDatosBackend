-- name: GetAllClients :many
SELECT * FROM clients;

-- name: GetClientById :one
SELECT * FROM clients WHERE client_id=?;

-- name: CreateClient :execresult
INSERT INTO clients (
    name,
    contact_phone,
    contact_email
)
VALUES (?, ?, ?);

-- name: UpdateClient :execresult
UPDATE clients
SET name = ?,
    contact_phone = ?,
    contact_email = ?
WHERE client_id=?;

-- name: DeleteClient :execresult
DELETE FROM clients WHERE client_id=?;
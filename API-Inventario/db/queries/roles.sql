-- name: GetAllRoles :many
SELECT * FROM roles;

-- name: GetRoleById :one
SELECT * FROM roles WHERE role_id=?;

-- name: CreateRole :execresult
INSERT INTO roles (nombre) VALUES (?);

-- name: UpdateRole :execresult
UPDATE roles SET nombre=? WHERE role_id=?;

-- name: DeleteRole :execresult
DELETE FROM roles WHERE role_id=?;
-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserById :one
SELECT * FROM users WHERE user_id=?;

-- name: GetAllUserRole :many
SELECT 
    u.user_id,
    u.name,
    u.password,
    u.email,
    r.role_id,
    r.name as role_name
FROM users u
LEFT JOIN roles r
ON u.role_id = r.role_id;



-- name: GetUserRoleById :one
SELECT 
    u.user_id,
    u.name,
    u.password,
    u.email,
    r.role_id,
    r.name as role_name
FROM users u
LEFT JOIN roles r
ON u.role_id = r.role_id
WHERE u.user_id = ?;

-- name: GetUserByEmail :one
SELECT 
    u.user_id,
    u.name,
    u.password,
    u.email,
    r.role_id,
    r.name as role_name
FROM users u
LEFT JOIN roles r
ON u.role_id = r.role_id
WHERE u.email = ?;

-- name: CreateUser :execresult
INSERT INTO users (role_id, name, password, email)
VALUES (?, ?, ?, ?);

-- name: UpdateUser :execresult
UPDATE users
SET role_id=?, name=?, password=?, email=?
WHERE user_id=?;

-- name: UpdateProfile :execresult
UPDATE users
SET
name = ?,
email = ?
WHERE user_id = ?;

-- name: UpdateUserRole :execresult
UPDATE users
SET role_id=?
WHERE user_id=?;

-- name: DeleteUser :execresult
DELETE FROM users WHERE user_id=?;
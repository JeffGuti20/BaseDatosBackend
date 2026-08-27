-- name: GetAllUsers :many
SELECT * FROM usuarios;

-- name: GetUserById :one
SELECT * FROM usuarios WHERE usuario_id=?;

-- name: GetUserByName :one
SELECT * FROM usuarios WHERE nombre=?;

-- name: GetAllUserRole :many
SELECT 
    u.usuario_id,
    u.nombre,
    u.contraseña,
    u.email,
    u.image,
    r.role_id,
    r.nombre as rol_nombre
FROM usuarios u
LEFT JOIN roles r 
ON u.rol_id = r.role_id;

-- name: GetUserRoleById :one
SELECT 
    u.usuario_id,
    u.nombre,
    u.contraseña,
    u.email,
    u.image,
    r.role_id,
    r.nombre as rol_nombre
FROM usuarios u
LEFT JOIN roles r 
ON u.rol_id = r.role_id
WHERE u.usuario_id = ?;

-- name: GetUserByEmail :one
SELECT 
    u.usuario_id,
    u.nombre,
    u.contraseña,
    u.email,
    u.image,
    r.role_id,
    r.nombre as rol_nombre
FROM usuarios u
LEFT JOIN roles r
ON u.rol_id = r.role_id
WHERE u.email = ?;

-- name: CreateUser :execresult
INSERT INTO usuarios (rol_id, nombre, contraseña, email, image)
VALUES (?, ?, ?, ?, ?);

-- name: UpdateUser :execresult
UPDATE usuarios
SET rol_id=?, nombre=?, contraseña=?, email=?, image=?
WHERE usuario_id=?;

-- name: UpdateUserRole :execresult
UPDATE usuarios
SET rol_id=?
WHERE usuario_id=?;

-- name: DeleteUser :execresult
DELETE FROM usuarios WHERE usuario_id=?;

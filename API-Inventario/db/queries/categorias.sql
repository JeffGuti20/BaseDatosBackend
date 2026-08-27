-- name: GetAllCategories :many
SELECT * FROM categorias;

-- name: GetCategoryById :one
SELECT * FROM categorias WHERE categoria_id=?;

-- name: CreateCategory :execresult
INSERT INTO categorias (nombre)
VALUES (?);

-- name: UpdateCategory :execresult
UPDATE categorias SET nombre=? WHERE categoria_id=?;

-- name: DeleteCategory :execresult
DELETE FROM categorias WHERE categoria_id=?;
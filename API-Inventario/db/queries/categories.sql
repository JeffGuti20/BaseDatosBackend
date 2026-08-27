-- name: GetAllCategories :many
SELECT * FROM categories;

-- name: GetCategoryById :one
SELECT * FROM categories WHERE category_id=?;

-- name: CreateCategory :execresult
INSERT INTO categories (name)
VALUES (?);

-- name: UpdateCategory :execresult
UPDATE categories SET name=? WHERE category_id=?;

-- name: DeleteCategory :execresult
DELETE FROM categories WHERE category_id=?;
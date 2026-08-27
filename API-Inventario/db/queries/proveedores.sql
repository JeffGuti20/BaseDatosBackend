-- name: GetAllSuppliers :many
SELECT * FROM proveedores;

-- name: GetSupplierById :one
SELECT * FROM proveedores WHERE proveedor_id=?;

-- name: CreateSupplier :execresult
INSERT INTO proveedores (nombre, telefono, email)
VALUES (?, ?, ?);

-- name: UpdateSupplier :execresult
UPDATE proveedores
SET nombre=?, telefono=?, email=?
WHERE proveedor_id=?;

-- name: DeleteSupplier :execresult
DELETE FROM proveedores WHERE proveedor_id=?;

-- name: GetAllMovements :many
SELECT
    m.movimiento_id,
    u.usuario_id,
    u.nombre as usuario_nombre,
    m.tipo,
    p.proveedor_id,
    p.nombre,
    p.telefono,
    p.email,
    m.fecha
FROM movimientos m
LEFT JOIN proveedores p
ON m.proveedor_id = p.proveedor_id
LEFT JOIN usuarios u
ON m.usuario_id = u.usuario_id;

-- name: GetMovementsById :one
SELECT
    m.movimiento_id,
    u.usuario_id,
    u.nombre as usuario_nombre,
    m.tipo,
    p.proveedor_id,
    p.nombre,
    p.telefono,
    p.email,
    m.fecha
FROM movimientos m
LEFT JOIN proveedores p
ON m.proveedor_id = p.proveedor_id
LEFT JOIN usuarios u
ON m.usuario_id = u.usuario_id
WHERE m.movimiento_id = ?;


-- name: CreateMovement :execresult
INSERT INTO movimientos (usuario_id, proveedor_id, tipo, fecha)
VALUES (?, ?, ?, now());

-- name: UpdateMovement :execresult
UPDATE movimientos
SET usuario_id=?, proveedor_id=?, tipo=?
WHERE movimiento_id=?;

-- name: DeleteMovement :execresult
DELETE FROM movimientos WHERE movimiento_id=?;

-- name: AddStock :exec
UPDATE insumos
SET cantidad = cantidad + ?
WHERE insumo_id = ?;

-- name: RemoveStock :exec
UPDATE insumos
SET cantidad = cantidad - ?
WHERE insumo_id = ?;
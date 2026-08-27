-- name: GetAllMovementDetails :many
SELECT 
    md.detalle_id,
    md.movimiento_id,
    md.insumo_id,
    md.cantidad,
    i.nombre,
    i.tipo
FROM movimiento_detalle md
LEFT JOIN insumos i
ON md.insumo_id = i.insumo_id;

-- name: GetDetailById :many
SELECT 
    md.detalle_id,
    md.movimiento_id,
    md.insumo_id,
    md.cantidad,
    i.nombre,
    i.tipo
FROM movimiento_detalle md
LEFT JOIN insumos i
ON md.insumo_id = i.insumo_id
WHERE md.movimiento_id=?;  


-- name: CreateMovementDetail :execresult
INSERT INTO movimiento_detalle (movimiento_id, insumo_id, cantidad)
VALUES (?, ?, ?);

-- name: UpdateMovementDetail :execresult
UPDATE movimiento_detalle
SET movimiento_id=?, insumo_id=?, cantidad=?
WHERE detalle_id=?;

-- name: DeleteMovementDetail :execresult
DELETE FROM movimiento_detalle WHERE movimiento_id=?;
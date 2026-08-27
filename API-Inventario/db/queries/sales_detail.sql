-- name: GetAllSaleDetails :many
SELECT
    sd.detail_id,
    sd.sale_id,
    sd.supply_id,
    sd.quantity,
    sd.unit_price,
    sd.discount,
    sd.subtotal,

    s.name as supply_name,
    s.type as supply_type
FROM sale_details sd
LEFT JOIN supplies s ON sd.supply_id = s.supply_id;

-- name: GetSaleDetailsBySaleId :many
SELECT
    sd.detail_id,
    sd.sale_id,
    sd.supply_id,
    sd.quantity,
    sd.unit_price,
    sd.discount,
    sd.subtotal,

    s.name as supply_name,
    s.type as supply_type
FROM sale_details sd
LEFT JOIN supplies s ON sd.supply_id = s.supply_id
WHERE sd.sale_id = ?;

-- name: CreateSaleDetail :execresult
INSERT INTO sale_details (
    sale_id,
    supply_id,
    quantity,
    unit_price,
    discount,
    subtotal
)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateSaleDetail :execresult
UPDATE sale_details
SET
    sale_id = ?,
    supply_id = ?,
    quantity = ?,
    unit_price = ?,
    discount = ?,
    subtotal = ?
WHERE detail_id = ?;

-- name: DeleteSaleDetail :execresult
DELETE FROM sale_details
WHERE sale_id = ?;
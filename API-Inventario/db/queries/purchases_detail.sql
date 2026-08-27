-- name: GetAllPurchaseDetails :many
SELECT
    pd.detail_id,
    pd.purchase_id,
    pd.supply_id,
    pd.quantity,
    pd.unit_price,
    pd.discount,
    pd.subtotal,

    s.name,
    s.type
FROM purchase_details pd
LEFT JOIN supplies s
ON pd.supply_id = s.supply_id;

-- name: GetPurchaseDetailsByPurchaseId :many
SELECT
    pd.detail_id,
    pd.purchase_id,
    pd.supply_id,
    pd.quantity,
    pd.unit_price,
    pd.discount,
    pd.subtotal,

    s.name,
    s.type
FROM purchase_details pd
LEFT JOIN supplies s
ON pd.supply_id = s.supply_id
WHERE pd.purchase_id=?;

-- name: CreatePurchaseDetail :execresult
INSERT INTO purchase_details (
    purchase_id,
    supply_id,
    quantity,
    unit_price,
    discount,
    subtotal
)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdatePurchaseDetail :execresult
UPDATE purchase_details
SET
    purchase_id=?,
    supply_id=?,
    quantity=?,
    unit_price=?,
    discount=?,
    subtotal=?
WHERE detail_id=?;

-- name: DeletePurchaseDetail :execresult
DELETE FROM purchase_details
WHERE purchase_id=?;
-- name: GetAllPurchases :many
SELECT
    p.purchase_id,
    p.date,
    p.subtotal,
    p.discount,
    p.IVA,
    p.total,

    u.user_id,
    u.name as user_name,

    s.supplier_id,
    s.company_name,
    s.company_phone,
    s.company_email,
    s.contact_name,
    s.contact_phone,
    s.contact_email
FROM purchases p
LEFT JOIN suppliers s
ON p.supplier_id = s.supplier_id
LEFT JOIN users u
ON p.user_id = u.user_id;

-- name: GetPurchaseById :one
SELECT
    p.purchase_id,
    p.date,
    p.subtotal,
    p.discount,
    p.IVA,
    p.total,

    u.user_id,
    u.name as user_name,

    s.supplier_id,
    s.company_name,
    s.company_phone,
    s.company_email,
    s.contact_name,
    s.contact_phone,
    s.contact_email
FROM purchases p
LEFT JOIN suppliers s
ON p.supplier_id = s.supplier_id
LEFT JOIN users u
ON p.user_id = u.user_id
WHERE p.purchase_id = ?;

-- name: CreatePurchase :execresult
INSERT INTO purchases (
    user_id,
    supplier_id,
    subtotal,
    discount,
    IVA,
    total
)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdatePurchase :execresult
UPDATE purchases
SET
    user_id=?,
    supplier_id=?,
    subtotal=?,
    discount=?,
    IVA=?,
    total=?
WHERE purchase_id=?;

-- name: DeletePurchase :execresult
DELETE FROM purchases
WHERE purchase_id=?;

-- name: AddStock :exec
UPDATE supplies
SET stock = stock + ?
WHERE supply_id = ?;

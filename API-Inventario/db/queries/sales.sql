-- name: GetAllSales :many
SELECT
    s.sale_id,
    s.date,
    s.subtotal,
    s.discount,
    s.IVA,
    s.total,

    u.user_id,
    u.name as user_name,

    c.client_id,
    c.name as client_name,
    c.contact_phone as client_phone,
    c.contact_email as client_email
FROM sales s
LEFT JOIN clients c ON s.client_id = c.client_id
LEFT JOIN users u ON s.user_id = u.user_id;

-- name: GetSaleById :one
SELECT
    s.sale_id,
    s.date,
    s.subtotal,
    s.discount,
    s.IVA,
    s.total,

    u.user_id,
    u.name as user_name,

    c.client_id,
    c.name as client_name,
    c.contact_phone as client_phone,
    c.contact_email as client_email
FROM sales s
LEFT JOIN clients c ON s.client_id = c.client_id
LEFT JOIN users u ON s.user_id = u.user_id
WHERE s.sale_id = ?;

-- name: CreateSale :execresult
INSERT INTO sales (
    user_id,
    client_id,
    subtotal,
    discount,
    IVA,
    total
)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateSale :execresult
UPDATE sales
SET
    user_id = ?,
    client_id = ?,
    subtotal = ?,
    discount = ?,
    IVA = ?,
    total = ?
WHERE sale_id = ?;

-- name: DeleteSale :execresult
DELETE FROM sales
WHERE sale_id = ?;


-- name: RemoveStock :exec
UPDATE supplies
SET stock = stock - ?
WHERE supply_id = ?;
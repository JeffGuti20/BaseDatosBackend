-- name: GetAllSuppliers :many
SELECT * FROM suppliers;

-- name: GetSupplierById :one
SELECT * FROM suppliers WHERE supplier_id=?;

-- name: CreateSupplier :execresult
INSERT INTO suppliers (
    company_name,
    company_phone,
    company_email,
    contact_name,
    contact_phone,
    contact_email
)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateSupplier :execresult
UPDATE suppliers
SET company_name=?,
    company_phone=?,
    company_email=?,
    contact_name=?,
    contact_phone=?,
    contact_email=?
WHERE supplier_id=?;

-- name: DeleteSupplier :execresult
DELETE FROM suppliers WHERE supplier_id=?;
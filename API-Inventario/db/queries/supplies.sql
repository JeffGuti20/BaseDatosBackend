-- name: GetAllSupplies :many
select s.supply_id,
       s.name,
       s.type,
       s.stock,
       s.category_id,
       s.price,
       s.image,
       c.name as category_name
from supplies s
left join categories c
on s.category_id = c.category_id;

-- name: GetSupplyById :one
select s.supply_id,
       s.name,
       s.type,
       s.stock,
       s.category_id,
       s.price,
       s.image,
       c.name as category_name
from supplies s
left join categories c
on s.category_id = c.category_id
where s.supply_id = ?;

-- name: CreateSupply :execresult
insert into supplies (
   name,
   type,
   stock,
   category_id,
   price,
   image
) values (?, ?, ?, ?, ?, ?);

-- name: UpdateSupply :execresult
update supplies
set name = ?,
    type = ?,
    stock = ?,
    category_id = ?,
    price = ?,
    image = ?
where supply_id = ?;

-- name: DeleteSupply :execresult
delete from supplies
where supply_id = ?;
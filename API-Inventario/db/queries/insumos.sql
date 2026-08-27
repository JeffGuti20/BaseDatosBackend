-- name: GetAllSupplies :many
select i.insumo_id,
       i.nombre,
       i.tipo,
       i.cantidad,
       i.image,
       i.categoria_id,
       c.nombre as categoria_nombre
from insumos i
left join categorias c
on i.categoria_id = c.categoria_id;

-- name: GetSupplyById :one
select i.insumo_id,
       i.nombre,
       i.tipo,
       i.cantidad,
       i.image,
       i.categoria_id,
       c.nombre as categoria_nombre
from insumos i
left join categorias c
on i.categoria_id = c.categoria_id
where i.insumo_id = ?;

-- name: CreateSupply :execresult
insert into insumos (
   nombre,
   tipo,
   cantidad,
   image,
   categoria_id
) values (?, ?, ?, ?, ?);

-- name: UpdateSupply :execresult
update insumos
set nombre = ?,
    tipo = ?,
    cantidad = ?,
    image = ?,
    categoria_id = ?
where insumo_id = ?;

-- name: DeleteSupply :execresult
delete from insumos
where insumo_id = ?;
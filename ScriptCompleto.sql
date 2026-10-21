/*
============================================================
EMI - INVENTARIO
SCRIPT SQL SERVER LIMPIO
============================================================
Base lógica de datos: inventariodb
Incluye:
- Tablas
- Claves primarias y foráneas
- Restricciones
- Valores por defecto
- Datos de prueba
- Vista
- Funciones
- Procedimientos almacenados

NOTA:
Los INSERT incluidos al final son DATOS DE PRUEBA agregados
para tener una base funcional. No corresponden necesariamente
a los datos reales de Coconutz.
============================================================
*/

-- ============================================================
-- 1. CREACIÓN DE BASE DE DATOS
-- ============================================================

IF DB_ID('inventariodb') IS NULL
BEGIN
    CREATE DATABASE inventariodb;
END;
GO

USE inventariodb;
GO

-- ============================================================
-- 2. TABLAS
-- ============================================================

CREATE TABLE categorias (
    idCategoria INT IDENTITY(1,1) NOT NULL,
    nombre VARCHAR(100) NOT NULL,
    CONSTRAINT PK_categorias PRIMARY KEY (idCategoria)
);
GO

CREATE TABLE roles (
    idRol INT IDENTITY(1,1) NOT NULL,
    nombre VARCHAR(300) NOT NULL,
    CONSTRAINT PK_roles PRIMARY KEY (idRol)
);
GO

CREATE TABLE clientes (
    idCliente INT IDENTITY(1,1) NOT NULL,
    nombre VARCHAR(300) NOT NULL,
    telefonoContacto INT NOT NULL,
    correoContacto VARCHAR(300) NOT NULL,
    CONSTRAINT PK_clientes PRIMARY KEY (idCliente)
);
GO

CREATE TABLE proveedores (
    idProveedor INT IDENTITY(1,1) NOT NULL,
    nombreEmpresa VARCHAR(300) NOT NULL,
    telefonoEmpresa INT NOT NULL,
    correoEmpresa VARCHAR(200) NOT NULL,
    nombreContacto VARCHAR(100) NULL,
    telefonoContacto INT NULL,
    correoContacto VARCHAR(300) NULL,
    CONSTRAINT PK_proveedores PRIMARY KEY (idProveedor)
);
GO

CREATE TABLE usuarios (
    idUsuario INT IDENTITY(1,1) NOT NULL,
    idRol INT NOT NULL,
    nombre VARCHAR(300) NOT NULL,
    contraseña VARCHAR(100) NOT NULL,
    correo VARCHAR(150) NULL,
    CONSTRAINT PK_usuarios PRIMARY KEY (idUsuario),
    CONSTRAINT FK_usuarios_roles
        FOREIGN KEY (idRol) REFERENCES roles(idRol)
);
GO

CREATE TABLE insumos (
    idInsumo INT IDENTITY(1,1) NOT NULL,
    nombre VARCHAR(300) NOT NULL,
    tipo VARCHAR(20) NOT NULL,
    stock INT NOT NULL CONSTRAINT DF_insumos_stock DEFAULT 0,
    idCategoria INT NOT NULL,
    precio DECIMAL(10,2) NOT NULL,
    imagen VARCHAR(300) NOT NULL,
    CONSTRAINT PK_insumos PRIMARY KEY (idInsumo),
    CONSTRAINT FK_insumos_categorias
        FOREIGN KEY (idCategoria) REFERENCES categorias(idCategoria),
    CONSTRAINT CK_insumos_tipo
        CHECK (tipo IN ('gastronomico', 'fisico'))
);
GO

CREATE TABLE compras (
    idCompra INT IDENTITY(1,1) NOT NULL,
    idUsuario INT NOT NULL,
    idProveedor INT NOT NULL,
    fecha DATETIME NULL CONSTRAINT DF_compras_fecha DEFAULT GETDATE(),
    subtotal DECIMAL(10,2) NULL,
    descuento DECIMAL(10,2) NULL,
    IVA DECIMAL(10,2) NOT NULL,
    total DECIMAL(10,2) NULL,
    CONSTRAINT PK_compras PRIMARY KEY (idCompra),
    CONSTRAINT FK_compras_usuarios
        FOREIGN KEY (idUsuario) REFERENCES usuarios(idUsuario),
    CONSTRAINT FK_compras_proveedores
        FOREIGN KEY (idProveedor) REFERENCES proveedores(idProveedor)
);
GO

CREATE TABLE ventas (
    idVenta INT IDENTITY(1,1) NOT NULL,
    idUsuario INT NOT NULL,
    idCliente INT NOT NULL,
    fecha DATETIME NULL CONSTRAINT DF_ventas_fecha DEFAULT GETDATE(),
    subtotal DECIMAL(10,2) NOT NULL,
    descuento DECIMAL(10,2) NULL,
    IVA DECIMAL(10,2) NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    CONSTRAINT PK_ventas PRIMARY KEY (idVenta),
    CONSTRAINT FK_ventas_usuarios
        FOREIGN KEY (idUsuario) REFERENCES usuarios(idUsuario),
    CONSTRAINT FK_ventas_clientes
        FOREIGN KEY (idCliente) REFERENCES clientes(idCliente)
);
GO

CREATE TABLE detallesCompra (
    idDetalle INT IDENTITY(1,1) NOT NULL,
    idCompra INT NOT NULL,
    idInsumo INT NOT NULL,
    cantidad INT NOT NULL CONSTRAINT DF_detallesCompra_cantidad DEFAULT 0,
    precioUnitario DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    descuento DECIMAL(10,2) NULL CONSTRAINT DF_detallesCompra_descuento DEFAULT 0,
    CONSTRAINT PK_detallesCompra PRIMARY KEY (idDetalle),
    CONSTRAINT FK_detallesCompra_compras
        FOREIGN KEY (idCompra) REFERENCES compras(idCompra),
    CONSTRAINT FK_detallesCompra_insumos
        FOREIGN KEY (idInsumo) REFERENCES insumos(idInsumo)
);
GO

CREATE TABLE detallesVenta (
    idDetalle INT IDENTITY(1,1) NOT NULL,
    idVenta INT NOT NULL,
    idInsumo INT NOT NULL,
    cantidad INT NOT NULL CONSTRAINT DF_detallesVenta_cantidad DEFAULT 0,
    precioUnitario DECIMAL(10,2) NOT NULL,
    descuento DECIMAL(10,2) NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    CONSTRAINT PK_detallesVenta PRIMARY KEY (idDetalle),
    CONSTRAINT FK_detallesVenta_ventas
        FOREIGN KEY (idVenta) REFERENCES ventas(idVenta),
    CONSTRAINT FK_detallesVenta_insumos
        FOREIGN KEY (idInsumo) REFERENCES insumos(idInsumo)
);
GO

-- ============================================================
-- 3. DATOS DE PRUEBA
-- ============================================================

INSERT INTO categorias (nombre)
VALUES
('Bebidas'),
('Alimentos'),
('Limpieza'),
('Empaque'),
('Equipo y utensilios');
GO

INSERT INTO roles (nombre)
VALUES
('Administrador'),
('Supervisor'),
('Operador');
GO

INSERT INTO clientes (nombre, telefonoContacto, correoContacto)
VALUES
('Cliente General', 88881111, 'cliente@ejemplo.com'),
('Hotel Coco Beach', 88882222, 'compras@cocobeach.com'),
('Restaurante El Coco', 88883333, 'admin@elcoco.com');
GO

INSERT INTO proveedores
(nombreEmpresa, telefonoEmpresa, correoEmpresa, nombreContacto, telefonoContacto, correoContacto)
VALUES
('Distribuidora Guanacaste', 26901234, 'ventas@distribuidorag.com',
 'Carlos Rodríguez', 87771111, 'carlos@distribuidorag.com'),
('Alimentos del Pacífico', 26905678, 'contacto@alimentospacifico.com',
 'María López', 87772222, 'maria@alimentospacifico.com'),
('Empaques CR', 22223333, 'ventas@empaquescr.com',
 'Daniel Vargas', 87773333, 'daniel@empaquescr.com');
GO

/*
Contraseñas de prueba.
En producción deben almacenarse como hashes.
Estos valores son únicamente datos de prueba.
*/
INSERT INTO usuarios (idRol, nombre, contraseña, correo)
VALUES
(1, 'Administrador EMI', 'hash_demo_admin', 'admin@emi.com'),
(2, 'Supervisor EMI', 'hash_demo_supervisor', 'supervisor@emi.com'),
(3, 'Operador EMI', 'hash_demo_operador', 'operador@emi.com');
GO

INSERT INTO insumos
(nombre, tipo, stock, idCategoria, precio, imagen)
VALUES
('Café en grano', 'gastronomico', 25, 2, 4500.00, 'cafe.jpg'),
('Leche entera', 'gastronomico', 18, 2, 1200.00, 'leche.jpg'),
('Jarabe de vainilla', 'gastronomico', 8, 1, 3500.00, 'vainilla.jpg'),
('Agua embotellada', 'gastronomico', 40, 1, 900.00, 'agua.jpg'),
('Vasos desechables', 'fisico', 120, 4, 35.00, 'vasos.jpg'),
('Servilletas', 'fisico', 200, 4, 15.00, 'servilletas.jpg'),
('Jabón líquido', 'fisico', 4, 3, 2500.00, 'jabon.jpg'),
('Guantes desechables', 'fisico', 0, 3, 1800.00, 'guantes.jpg'),
('Cucharas plásticas', 'fisico', 70, 4, 20.00, 'cucharas.jpg'),
('Tazas reutilizables', 'fisico', 6, 5, 2500.00, 'tazas.jpg');
GO

INSERT INTO compras
(idUsuario, idProveedor, fecha, subtotal, descuento, IVA, total)
VALUES
(1, 1, '2026-09-25 09:30:00', 95000.00, 5000.00, 13.00, 101700.00),
(2, 2, '2026-09-27 10:15:00', 60000.00, 0.00, 13.00, 67800.00),
(1, 3, '2026-09-29 14:00:00', 25000.00, 1000.00, 13.00, 27120.00);
GO

INSERT INTO detallesCompra
(idCompra, idInsumo, cantidad, precioUnitario, subtotal, descuento)
VALUES
(1, 1, 10, 4500.00, 45000.00, 2500.00),
(1, 3, 5, 3500.00, 17500.00, 1000.00),
(1, 5, 500, 35.00, 17500.00, 1000.00),
(2, 2, 30, 1200.00, 36000.00, 0.00),
(2, 4, 20, 900.00, 18000.00, 0.00),
(3, 6, 1000, 15.00, 15000.00, 500.00),
(3, 9, 500, 20.00, 10000.00, 500.00);
GO

INSERT INTO ventas
(idUsuario, idCliente, fecha, subtotal, descuento, IVA, total)
VALUES
(3, 1, '2026-09-28 11:20:00', 12000.00, 500.00, 13.00, 13005.00),
(2, 2, '2026-09-29 15:40:00', 22500.00, 1000.00, 13.00, 24295.00),
(3, 3, '2026-09-30 10:10:00', 7000.00, 0.00, 13.00, 7910.00);
GO

INSERT INTO detallesVenta
(idVenta, idInsumo, cantidad, precioUnitario, descuento, subtotal)
VALUES
(1, 1, 2, 4500.00, 500.00, 9000.00),
(1, 4, 4, 750.00, 0.00, 3000.00),
(2, 3, 3, 3500.00, 500.00, 10500.00),
(2, 5, 300, 40.00, 500.00, 12000.00),
(3, 6, 100, 15.00, 0.00, 1500.00),
(3, 9, 250, 20.00, 0.00, 5000.00);
GO

-- ============================================================
-- 4. VISTA DE MOVIMIENTOS CON DESCUENTO
-- ============================================================

CREATE VIEW vMovimientosDescuento
AS
    SELECT
        idCompra AS idMovimiento,
        subtotal,
        ISNULL(descuento, 0) AS descuento,
        'Compras' AS tipo
    FROM detallesCompra

    UNION ALL

    SELECT
        idVenta AS idMovimiento,
        subtotal,
        ISNULL(descuento, 0) AS descuento,
        'Ventas' AS tipo
    FROM detallesVenta;
GO

-- ============================================================
-- 5. FUNCIONES ESCALARES
-- ============================================================

CREATE FUNCTION fnCalcularDescuento
(
    @id INT,
    @tabla VARCHAR(10)
)
RETURNS DECIMAL(10,2)
AS
BEGIN
    DECLARE @resultado DECIMAL(10,2);

    SELECT @resultado =
        SUM(subtotal - (subtotal * ISNULL(descuento, 0) / 100))
    FROM vMovimientosDescuento
    WHERE idMovimiento = @id
      AND tipo = @tabla;

    RETURN ISNULL(@resultado, 0);
END;
GO

CREATE FUNCTION fnCalcularTotal
(
    @subtotal DECIMAL(10,2),
    @descuento DECIMAL(10,2),
    @iva DECIMAL(10,2)
)
RETURNS DECIMAL(10,2)
AS
BEGIN
    RETURN
        (@subtotal - @descuento)
        + ((@subtotal - @descuento) * @iva / 100);
END;
GO

-- ============================================================
-- 6. FUNCIONES CON VALORES DE TABLA
-- ============================================================

CREATE FUNCTION fnMovimientosInsumo
(
    @IdInsumo INT
)
RETURNS TABLE
AS
RETURN
(
    SELECT
        'Compra' AS tipoMovimiento,
        dc.idCompra AS idMovimiento,
        c.fecha,
        dc.cantidad,
        dc.precioUnitario,
        dc.subtotal,
        ISNULL(dc.descuento, 0) AS descuento
    FROM detallesCompra dc
    INNER JOIN compras c
        ON c.idCompra = dc.idCompra
    WHERE dc.idInsumo = @IdInsumo

    UNION ALL

    SELECT
        'Venta' AS tipoMovimiento,
        dv.idVenta AS idMovimiento,
        v.fecha,
        dv.cantidad,
        dv.precioUnitario,
        dv.subtotal,
        ISNULL(dv.descuento, 0) AS descuento
    FROM detallesVenta dv
    INNER JOIN ventas v
        ON v.idVenta = dv.idVenta
    WHERE dv.idInsumo = @IdInsumo
);
GO

CREATE FUNCTION fnVentasCliente
(
    @IdCliente INT
)
RETURNS TABLE
AS
RETURN
(
    SELECT
        v.idVenta,
        v.fecha,
        dv.idInsumo,
        dv.cantidad,
        dv.precioUnitario,
        dv.subtotal,
        ISNULL(dv.descuento, 0) AS descuento
    FROM ventas v
    INNER JOIN detallesVenta dv
        ON dv.idVenta = v.idVenta
    WHERE v.idCliente = @IdCliente
);
GO

CREATE FUNCTION fnInsumosAgotados()
RETURNS TABLE
AS
RETURN
(
    SELECT
        i.idInsumo,
        i.nombre,
        i.tipo,
        i.stock,
        i.idCategoria,
        c.nombre AS categoria,
        i.precio,
        i.imagen
    FROM insumos i
    INNER JOIN categorias c
        ON c.idCategoria = i.idCategoria
    WHERE i.stock = 0
);
GO

CREATE FUNCTION fnInsumosStockBajo()
RETURNS TABLE
AS
RETURN
(
    SELECT
        i.idInsumo,
        i.nombre,
        i.tipo,
        i.stock,
        i.idCategoria,
        c.nombre AS categoria,
        i.precio,
        i.imagen
    FROM insumos i
    INNER JOIN categorias c
        ON c.idCategoria = i.idCategoria
    WHERE i.stock <= 5
);
GO

-- ============================================================
-- 7. PROCEDIMIENTOS - CATEGORÍAS
-- ============================================================

CREATE PROCEDURE paCrearCategoria
    @Nombre VARCHAR(100)
AS
BEGIN
    SET NOCOUNT ON;

    SET @Nombre = LTRIM(RTRIM(@Nombre));

    IF @Nombre = ''
        RAISERROR('El nombre de la categoría no puede estar vacío.', 16, 1);

    IF EXISTS (
        SELECT 1
        FROM categorias
        WHERE LOWER(nombre) = LOWER(@Nombre)
    )
        RAISERROR('Ya existe una categoría con ese nombre.', 16, 1);

    INSERT INTO categorias (nombre)
    VALUES (@Nombre);

    SELECT SCOPE_IDENTITY() AS idCategoria;
END;
GO

CREATE PROCEDURE paActualizarCategoria
    @IdCategoria INT,
    @Nombre VARCHAR(100)
AS
BEGIN
    SET NOCOUNT ON;

    SET @Nombre = LTRIM(RTRIM(@Nombre));

    IF @IdCategoria IS NULL OR @IdCategoria <= 0
        RAISERROR('El ID de la categoría no es válido.', 16, 1);

    IF @Nombre = ''
        RAISERROR('El nombre de la categoría no puede estar vacío.', 16, 1);

    IF NOT EXISTS (
        SELECT 1 FROM categorias WHERE idCategoria = @IdCategoria
    )
        RAISERROR('La categoría no existe.', 16, 1);

    IF EXISTS (
        SELECT 1
        FROM categorias
        WHERE LOWER(nombre) = LOWER(@Nombre)
          AND idCategoria <> @IdCategoria
    )
        RAISERROR('Ya existe otra categoría con ese nombre.', 16, 1);

    UPDATE categorias
    SET nombre = @Nombre
    WHERE idCategoria = @IdCategoria;

    SELECT *
    FROM categorias
    WHERE idCategoria = @IdCategoria;
END;
GO

CREATE PROCEDURE paEliminarCategoria
    @IdCategoria INT
AS
BEGIN
    SET NOCOUNT ON;

    IF @IdCategoria IS NULL OR @IdCategoria <= 0
        RAISERROR('El ID de la categoría no es válido.', 16, 1);

    IF NOT EXISTS (
        SELECT 1 FROM categorias WHERE idCategoria = @IdCategoria
    )
        RAISERROR('La categoría no existe.', 16, 1);

    IF EXISTS (
        SELECT 1
        FROM insumos
        WHERE idCategoria = @IdCategoria
    )
        RAISERROR('No se puede eliminar la categoría porque tiene insumos asociados.', 16, 1);

    DELETE FROM categorias
    WHERE idCategoria = @IdCategoria;
END;
GO

CREATE PROCEDURE paObtenerCategoria
    @IdCategoria INT
AS
BEGIN
    SET NOCOUNT ON;

    IF @IdCategoria IS NULL OR @IdCategoria <= 0
        RAISERROR('El ID de la categoría no es válido.', 16, 1);

    IF NOT EXISTS (
        SELECT 1 FROM categorias WHERE idCategoria = @IdCategoria
    )
        RAISERROR('La categoría no existe.', 16, 1);

    SELECT *
    FROM categorias
    WHERE idCategoria = @IdCategoria;
END;
GO

CREATE PROCEDURE paObtenerCategorias
AS
BEGIN
    SET NOCOUNT ON;

    SELECT *
    FROM categorias
    ORDER BY idCategoria;
END;
GO

-- ============================================================
-- 8. PROCEDIMIENTOS - INSUMOS
-- ============================================================

CREATE PROCEDURE paCrearInsumo
    @Nombre VARCHAR(300),
    @Tipo VARCHAR(20),
    @Stock INT,
    @IdCategoria INT,
    @Precio DECIMAL(10,2),
    @Imagen VARCHAR(300)
AS
BEGIN
    SET NOCOUNT ON;

    SET @Nombre = LTRIM(RTRIM(@Nombre));
    SET @Tipo = LTRIM(RTRIM(@Tipo));
    SET @Imagen = LTRIM(RTRIM(@Imagen));

    IF @Nombre = ''
        RAISERROR('El nombre del insumo no puede estar vacío.', 16, 1);

    IF @Tipo NOT IN ('fisico', 'gastronomico')
        RAISERROR('El tipo de insumo no es válido.', 16, 1);

    IF @Stock < 0
        RAISERROR('El stock no puede ser negativo.', 16, 1);

    IF @IdCategoria IS NULL OR @IdCategoria <= 0
        RAISERROR('La categoría no es válida.', 16, 1);

    IF NOT EXISTS (
        SELECT 1 FROM categorias WHERE idCategoria = @IdCategoria
    )
        RAISERROR('La categoría no existe.', 16, 1);

    IF @Precio < 0
        RAISERROR('El precio no puede ser negativo.', 16, 1);

    IF @Imagen = ''
        RAISERROR('La imagen no puede estar vacía.', 16, 1);

    IF EXISTS (
        SELECT 1
        FROM insumos
        WHERE LOWER(nombre) = LOWER(@Nombre)
    )
        RAISERROR('Ya existe un insumo con ese nombre.', 16, 1);

    INSERT INTO insumos
    (nombre, tipo, stock, idCategoria, precio, imagen)
    VALUES
    (@Nombre, @Tipo, @Stock, @IdCategoria, @Precio, @Imagen);

    SELECT SCOPE_IDENTITY() AS idInsumo;
END;
GO

CREATE PROCEDURE paActualizarInsumo
    @IdInsumo INT,
    @Nombre VARCHAR(300),
    @Tipo VARCHAR(20),
    @Stock INT,
    @IdCategoria INT,
    @Precio DECIMAL(10,2),
    @Imagen VARCHAR(300)
AS
BEGIN
    SET NOCOUNT ON;

    SET @Nombre = LTRIM(RTRIM(@Nombre));
    SET @Tipo = LTRIM(RTRIM(@Tipo));
    SET @Imagen = LTRIM(RTRIM(@Imagen));

    IF @IdInsumo IS NULL OR @IdInsumo <= 0
        RAISERROR('El ID del insumo no es válido.', 16, 1);

    IF NOT EXISTS (
        SELECT 1 FROM insumos WHERE idInsumo = @IdInsumo
    )
        RAISERROR('El insumo no existe.', 16, 1);

    IF @Nombre = ''
        RAISERROR('El nombre del insumo no puede estar vacío.', 16, 1);

    IF @Tipo NOT IN ('fisico', 'gastronomico')
        RAISERROR('El tipo de insumo no es válido.', 16, 1);

    IF @Stock < 0
        RAISERROR('El stock no puede ser negativo.', 16, 1);

    IF @IdCategoria IS NULL OR @IdCategoria <= 0
        RAISERROR('La categoría no es válida.', 16, 1);

    IF NOT EXISTS (
        SELECT 1 FROM categorias WHERE idCategoria = @IdCategoria
    )
        RAISERROR('La categoría no existe.', 16, 1);

    IF @Precio < 0
        RAISERROR('El precio no puede ser negativo.', 16, 1);

    IF @Imagen = ''
        RAISERROR('La imagen no puede estar vacía.', 16, 1);

    IF EXISTS (
        SELECT 1
        FROM insumos
        WHERE LOWER(nombre) = LOWER(@Nombre)
          AND idInsumo <> @IdInsumo
    )
        RAISERROR('Ya existe otro insumo con ese nombre.', 16, 1);

    UPDATE insumos
    SET
        nombre = @Nombre,
        tipo = @Tipo,
        stock = @Stock,
        idCategoria = @IdCategoria,
        precio = @Precio,
        imagen = @Imagen
    WHERE idInsumo = @IdInsumo;

    SELECT
        i.*,
        c.nombre AS categoria
    FROM insumos i
    INNER JOIN categorias c
        ON c.idCategoria = i.idCategoria
    WHERE i.idInsumo = @IdInsumo;
END;
GO

CREATE PROCEDURE paEliminarInsumo
    @IdInsumo INT
AS
BEGIN
    SET NOCOUNT ON;

    IF @IdInsumo IS NULL OR @IdInsumo <= 0
        RAISERROR('El ID del insumo no es válido.', 16, 1);

    IF NOT EXISTS (
        SELECT 1 FROM insumos WHERE idInsumo = @IdInsumo
    )
        RAISERROR('El insumo no existe.', 16, 1);

    IF EXISTS (
        SELECT 1
        FROM detallesCompra
        WHERE idInsumo = @IdInsumo
    )
        RAISERROR('No se puede eliminar el insumo porque tiene compras asociadas.', 16, 1);

    IF EXISTS (
        SELECT 1
        FROM detallesVenta
        WHERE idInsumo = @IdInsumo
    )
        RAISERROR('No se puede eliminar el insumo porque tiene ventas asociadas.', 16, 1);

    DELETE FROM insumos
    WHERE idInsumo = @IdInsumo;
END;
GO

CREATE PROCEDURE paObtenerInsumo
    @IdInsumo INT
AS
BEGIN
    SET NOCOUNT ON;

    IF @IdInsumo IS NULL OR @IdInsumo <= 0
        RAISERROR('El ID del insumo no es válido.', 16, 1);

    IF NOT EXISTS (
        SELECT 1 FROM insumos WHERE idInsumo = @IdInsumo
    )
        RAISERROR('El insumo no existe.', 16, 1);

    SELECT
        i.*,
        c.nombre AS categoria
    FROM insumos i
    INNER JOIN categorias c
        ON c.idCategoria = i.idCategoria
    WHERE i.idInsumo = @IdInsumo;
END;
GO

CREATE PROCEDURE paObtenerInsumos
AS
BEGIN
    SET NOCOUNT ON;

    SELECT
        i.*,
        c.nombre AS categoria
    FROM insumos i
    INNER JOIN categorias c
        ON c.idCategoria = i.idCategoria
    ORDER BY i.idInsumo;
END;
GO

-- ============================================================
-- 9. PROCEDIMIENTOS - ROLES
-- ============================================================

CREATE PROCEDURE paCrearRol
    @Nombre VARCHAR(300)
AS
BEGIN
    SET NOCOUNT ON;

    SET @Nombre = LTRIM(RTRIM(@Nombre));

    IF @Nombre = ''
        RAISERROR('El nombre del rol no puede estar vacío.', 16, 1);

    IF EXISTS (
        SELECT 1
        FROM roles
        WHERE LOWER(nombre) = LOWER(@Nombre)
    )
        RAISERROR('Ya existe un rol con ese nombre.', 16, 1);

    INSERT INTO roles (nombre)
    VALUES (@Nombre);

    SELECT SCOPE_IDENTITY() AS idRole;
END;
GO

CREATE PROCEDURE paActualizarRol
    @IdRol INT,
    @Nombre VARCHAR(300)
AS
BEGIN
    SET NOCOUNT ON;

    SET @Nombre = LTRIM(RTRIM(@Nombre));

    IF @IdRol IS NULL OR @IdRol <= 0
        RAISERROR('El ID del rol no es válido.', 16, 1);

    IF @Nombre = ''
        RAISERROR('El nombre del rol no puede estar vacío.', 16, 1);

    IF NOT EXISTS (
        SELECT 1 FROM roles WHERE idRol = @IdRol
    )
        RAISERROR('El rol no existe.', 16, 1);

    IF EXISTS (
        SELECT 1
        FROM roles
        WHERE LOWER(nombre) = LOWER(@Nombre)
          AND idRol <> @IdRol
    )
        RAISERROR('Ya existe otro rol con ese nombre.', 16, 1);

    UPDATE roles
    SET nombre = @Nombre
    WHERE idRol = @IdRol;

    SELECT *
    FROM roles
    WHERE idRol = @IdRol;
END;
GO

CREATE PROCEDURE paEliminarRol
    @IdRol INT
AS
BEGIN
    SET NOCOUNT ON;

    IF @IdRol IS NULL OR @IdRol <= 0
        RAISERROR('El ID del rol no es válido.', 16, 1);

    IF NOT EXISTS (
        SELECT 1 FROM roles WHERE idRol = @IdRol
    )
        RAISERROR('El rol no existe.', 16, 1);

    IF EXISTS (
        SELECT 1
        FROM usuarios
        WHERE idRol = @IdRol
    )
        RAISERROR('No se puede eliminar el rol porque tiene usuarios asociados.', 16, 1);

    DELETE FROM roles
    WHERE idRol = @IdRol;
END;
GO

CREATE PROCEDURE paObtenerRol
    @IdRol INT
AS
BEGIN
    SET NOCOUNT ON;

    IF @IdRol IS NULL OR @IdRol <= 0
        RAISERROR('El ID del rol no es válido.', 16, 1);

    IF NOT EXISTS (
        SELECT 1 FROM roles WHERE idRol = @IdRol
    )
        RAISERROR('El rol no existe.', 16, 1);

    SELECT *
    FROM roles
    WHERE idRol = @IdRol;
END;
GO

CREATE PROCEDURE paObtenerRoles
AS
BEGIN
    SET NOCOUNT ON;

    SELECT *
    FROM roles
    ORDER BY idRol;
END;
GO

-- ============================================================
-- 10. PROCEDIMIENTO - INICIO DE SESIÓN
-- ============================================================

CREATE PROCEDURE paIniciarSesion
    @Correo VARCHAR(150)
AS
BEGIN
    SET NOCOUNT ON;

    SET @Correo = LTRIM(RTRIM(@Correo));

    IF @Correo = ''
        RAISERROR('El correo no puede estar vacío.', 16, 1);

    IF @Correo NOT LIKE '%_@_%._%'
        RAISERROR('El formato del correo no es válido.', 16, 1);

    IF NOT EXISTS (
        SELECT 1
        FROM usuarios
        WHERE correo = @Correo
    )
        RAISERROR('No existe un usuario con ese correo.', 16, 1);

    SELECT
        u.idUsuario,
        u.idRol,
        u.nombre,
        u.contraseña,
        u.correo,
        r.nombre AS rol
    FROM usuarios u
    INNER JOIN roles r
        ON r.idRol = u.idRol
    WHERE u.correo = @Correo;
END;
GO

-- ============================================================
-- 11. PRUEBAS RÁPIDAS
-- ============================================================

-- SELECT * FROM categorias;
-- SELECT * FROM roles;
-- SELECT * FROM usuarios;
-- SELECT * FROM insumos;
-- SELECT * FROM proveedores;
-- SELECT * FROM clientes;
-- SELECT * FROM compras;
-- SELECT * FROM detallesCompra;
-- SELECT * FROM ventas;
-- SELECT * FROM detallesVenta;

-- Ejemplos de funciones:
-- SELECT dbo.fnCalcularDescuento(1, 'Compras') AS descuento;
-- SELECT dbo.fnCalcularTotal(10000, 500, 13) AS total;
-- SELECT * FROM dbo.fnMovimientosInsumo(1);
-- SELECT * FROM dbo.fnVentasCliente(1);
-- SELECT * FROM dbo.fnInsumosAgotados();
-- SELECT * FROM dbo.fnInsumosStockBajo();

-- Ejemplos de procedimientos:
-- EXEC paObtenerCategorias;
-- EXEC paObtenerInsumos;
-- EXEC paObtenerRoles;
-- EXEC paIniciarSesion 'admin@emi.com';

-- ============================================================
-- 12. AUDITORIA DE USUARIOS, CLIENTES E INSUMOS
-- ============================================================
-- Los triggers registran INSERT, UPDATE y DELETE.
-- La identidad del usuario autenticado por EMI se obtiene de
-- SESSION_CONTEXT('idUsuario'). La API Go debe establecer ese
-- contexto en la misma conexión SQL que realiza cada operación
-- auditada y limpiarlo antes de devolver la conexión al pool.
-- Ejemplo de configuración desde la API (misma conexión):
-- EXEC sys.sp_set_session_context @key=N'idUsuario', @value=@idUsuario;
-- No se almacena la contraseña ni su hash en el historial.

USE inventariodb;
GO

/* =============================================================
   1. TABLAS DE AUDITORIA
   ============================================================= */

IF OBJECT_ID(N'dbo.auditoriaUsuarios', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.auditoriaUsuarios
    (
        idAuditoria BIGINT IDENTITY(1,1) NOT NULL
            CONSTRAINT PK_auditoriaUsuarios PRIMARY KEY,
        tablaAfectada SYSNAME NOT NULL
            CONSTRAINT DF_auditoriaUsuarios_tabla DEFAULT (N'usuarios'),
        operacion VARCHAR(10) NOT NULL,
        idRegistro INT NOT NULL,
        fechaHora DATETIME2(3) NOT NULL
            CONSTRAINT DF_auditoriaUsuarios_fecha DEFAULT (SYSDATETIME()),
        loginSQL NVARCHAR(128) NULL,
        idUsuarioAplicacion INT NULL,
        valoresAnteriores NVARCHAR(MAX) NULL,
        valoresNuevos NVARCHAR(MAX) NULL,
        contrasenaCambiada BIT NOT NULL
            CONSTRAINT DF_auditoriaUsuarios_contrasena DEFAULT (0),
        CONSTRAINT CK_auditoriaUsuarios_operacion
            CHECK (operacion IN ('INSERT', 'UPDATE', 'DELETE'))
    );
END;
GO

IF OBJECT_ID(N'dbo.auditoriaClientes', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.auditoriaClientes
    (
        idAuditoria BIGINT IDENTITY(1,1) NOT NULL
            CONSTRAINT PK_auditoriaClientes PRIMARY KEY,
        tablaAfectada SYSNAME NOT NULL
            CONSTRAINT DF_auditoriaClientes_tabla DEFAULT (N'clientes'),
        operacion VARCHAR(10) NOT NULL,
        idRegistro INT NOT NULL,
        fechaHora DATETIME2(3) NOT NULL
            CONSTRAINT DF_auditoriaClientes_fecha DEFAULT (SYSDATETIME()),
        loginSQL NVARCHAR(128) NULL,
        idUsuarioAplicacion INT NULL,
        valoresAnteriores NVARCHAR(MAX) NULL,
        valoresNuevos NVARCHAR(MAX) NULL,
        CONSTRAINT CK_auditoriaClientes_operacion
            CHECK (operacion IN ('INSERT', 'UPDATE', 'DELETE'))
    );
END;
GO

IF OBJECT_ID(N'dbo.auditoriaInsumos', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.auditoriaInsumos
    (
        idAuditoria BIGINT IDENTITY(1,1) NOT NULL
            CONSTRAINT PK_auditoriaInsumos PRIMARY KEY,
        tablaAfectada SYSNAME NOT NULL
            CONSTRAINT DF_auditoriaInsumos_tabla DEFAULT (N'insumos'),
        operacion VARCHAR(10) NOT NULL,
        idRegistro INT NOT NULL,
        fechaHora DATETIME2(3) NOT NULL
            CONSTRAINT DF_auditoriaInsumos_fecha DEFAULT (SYSDATETIME()),
        loginSQL NVARCHAR(128) NULL,
        idUsuarioAplicacion INT NULL,
        valoresAnteriores NVARCHAR(MAX) NULL,
        valoresNuevos NVARCHAR(MAX) NULL,
        CONSTRAINT CK_auditoriaInsumos_operacion
            CHECK (operacion IN ('INSERT', 'UPDATE', 'DELETE'))
    );
END;
GO

/* =============================================================
   2. TRIGGER DE USUARIOS
   Se excluye la columna [contraseña] de los JSON por seguridad.
   Si cambia, se registra solo la marca contrasenaCambiada = 1.
   ============================================================= */
CREATE OR ALTER TRIGGER dbo.trg_AuditoriaUsuarios
ON dbo.usuarios
AFTER INSERT, UPDATE, DELETE
AS
BEGIN
    SET NOCOUNT ON;

    INSERT INTO dbo.auditoriaUsuarios
    (
        tablaAfectada, operacion, idRegistro, fechaHora,
        loginSQL, idUsuarioAplicacion,
        valoresAnteriores, valoresNuevos, contrasenaCambiada
    )
    SELECT
        N'usuarios',
        CASE
            WHEN d.idUsuario IS NULL THEN 'INSERT'
            WHEN i.idUsuario IS NULL THEN 'DELETE'
            ELSE 'UPDATE'
        END,
        COALESCE(i.idUsuario, d.idUsuario),
        SYSDATETIME(),
        ORIGINAL_LOGIN(),
        TRY_CONVERT(INT, SESSION_CONTEXT(N'idUsuario')),
        CASE WHEN d.idUsuario IS NULL THEN NULL ELSE
            (SELECT d.idUsuario AS idUsuario,
                    d.idRol AS idRol,
                    d.nombre AS nombre,
                    d.correo AS correo
             FOR JSON PATH, INCLUDE_NULL_VALUES, WITHOUT_ARRAY_WRAPPER)
        END,
        CASE WHEN i.idUsuario IS NULL THEN NULL ELSE
            (SELECT i.idUsuario AS idUsuario,
                    i.idRol AS idRol,
                    i.nombre AS nombre,
                    i.correo AS correo
             FOR JSON PATH, INCLUDE_NULL_VALUES, WITHOUT_ARRAY_WRAPPER)
        END,
        CASE
            WHEN i.idUsuario IS NOT NULL
             AND d.idUsuario IS NOT NULL
             AND ISNULL(i.[contraseña], '') <> ISNULL(d.[contraseña], '')
            THEN 1
            ELSE 0
        END
    FROM inserted AS i
    FULL OUTER JOIN deleted AS d
        ON i.idUsuario = d.idUsuario;
END;
GO

/* =============================================================
   3. TRIGGER DE CLIENTES
   ============================================================= */
CREATE OR ALTER TRIGGER dbo.trg_AuditoriaClientes
ON dbo.clientes
AFTER INSERT, UPDATE, DELETE
AS
BEGIN
    SET NOCOUNT ON;

    INSERT INTO dbo.auditoriaClientes
    (
        tablaAfectada, operacion, idRegistro, fechaHora,
        loginSQL, idUsuarioAplicacion,
        valoresAnteriores, valoresNuevos
    )
    SELECT
        N'clientes',
        CASE
            WHEN d.idCliente IS NULL THEN 'INSERT'
            WHEN i.idCliente IS NULL THEN 'DELETE'
            ELSE 'UPDATE'
        END,
        COALESCE(i.idCliente, d.idCliente),
        SYSDATETIME(),
        ORIGINAL_LOGIN(),
        TRY_CONVERT(INT, SESSION_CONTEXT(N'idUsuario')),
        CASE WHEN d.idCliente IS NULL THEN NULL ELSE
            (SELECT d.idCliente AS idCliente,
                    d.nombre AS nombre,
                    d.telefonoContacto AS telefonoContacto,
                    d.correoContacto AS correoContacto
             FOR JSON PATH, INCLUDE_NULL_VALUES, WITHOUT_ARRAY_WRAPPER)
        END,
        CASE WHEN i.idCliente IS NULL THEN NULL ELSE
            (SELECT i.idCliente AS idCliente,
                    i.nombre AS nombre,
                    i.telefonoContacto AS telefonoContacto,
                    i.correoContacto AS correoContacto
             FOR JSON PATH, INCLUDE_NULL_VALUES, WITHOUT_ARRAY_WRAPPER)
        END
    FROM inserted AS i
    FULL OUTER JOIN deleted AS d
        ON i.idCliente = d.idCliente;
END;
GO

/* =============================================================
   4. TRIGGER DE INSUMOS
   ============================================================= */
CREATE OR ALTER TRIGGER dbo.trg_AuditoriaInsumos
ON dbo.insumos
AFTER INSERT, UPDATE, DELETE
AS
BEGIN
    SET NOCOUNT ON;

    INSERT INTO dbo.auditoriaInsumos
    (
        tablaAfectada, operacion, idRegistro, fechaHora,
        loginSQL, idUsuarioAplicacion,
        valoresAnteriores, valoresNuevos
    )
    SELECT
        N'insumos',
        CASE
            WHEN d.idInsumo IS NULL THEN 'INSERT'
            WHEN i.idInsumo IS NULL THEN 'DELETE'
            ELSE 'UPDATE'
        END,
        COALESCE(i.idInsumo, d.idInsumo),
        SYSDATETIME(),
        ORIGINAL_LOGIN(),
        TRY_CONVERT(INT, SESSION_CONTEXT(N'idUsuario')),
        CASE WHEN d.idInsumo IS NULL THEN NULL ELSE
            (SELECT d.idInsumo AS idInsumo,
                    d.nombre AS nombre,
                    d.tipo AS tipo,
                    d.stock AS stock,
                    d.idCategoria AS idCategoria,
                    d.precio AS precio,
                    d.imagen AS imagen
             FOR JSON PATH, INCLUDE_NULL_VALUES, WITHOUT_ARRAY_WRAPPER)
        END,
        CASE WHEN i.idInsumo IS NULL THEN NULL ELSE
            (SELECT i.idInsumo AS idInsumo,
                    i.nombre AS nombre,
                    i.tipo AS tipo,
                    i.stock AS stock,
                    i.idCategoria AS idCategoria,
                    i.precio AS precio,
                    i.imagen AS imagen
             FOR JSON PATH, INCLUDE_NULL_VALUES, WITHOUT_ARRAY_WRAPPER)
        END
    FROM inserted AS i
    FULL OUTER JOIN deleted AS d
        ON i.idInsumo = d.idInsumo;
END;
GO

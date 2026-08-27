CREATE TABLE categorias (
  categoria_id INT NOT NULL AUTO_INCREMENT,
  nombre VARCHAR(100) NOT NULL,
  PRIMARY KEY (categoria_id)
);

CREATE TABLE insumos (
  insumo_id INT NOT NULL AUTO_INCREMENT,
  nombre VARCHAR(300) NOT NULL,
  tipo ENUM('fisico','gastronomico') NOT NULL,
  cantidad INT NOT NULL DEFAULT 0,
  categoria_id INT NOT NULL DEFAULT 0,
  image VARCHAR(300) NOT NULL,
  PRIMARY KEY (insumo_id),
  KEY fk_productos_categoria (categoria_id),
  CONSTRAINT fk_productos_categoria
    FOREIGN KEY (categoria_id)
    REFERENCES categorias (categoria_id)
    ON DELETE SET NULL
);

CREATE TABLE proveedores (
  proveedor_id INT NOT NULL AUTO_INCREMENT,
  nombre VARCHAR(300) NOT NULL,
  telefono INT DEFAULT NULL,
  email VARCHAR(300) DEFAULT NULL,
  PRIMARY KEY (proveedor_id)
);

CREATE TABLE roles (
  role_id INT NOT NULL AUTO_INCREMENT,
  nombre VARCHAR(300) NOT NULL,
  PRIMARY KEY (role_id)
);

CREATE TABLE usuarios (
  usuario_id INT NOT NULL AUTO_INCREMENT,
  rol_id INT NOT NULL,
  nombre VARCHAR(300) NOT NULL,
  contraseña VARCHAR(100) NOT NULL,
  email VARCHAR(150) DEFAULT NULL,
  image VARCHAR(300) NOT NULL,
  PRIMARY KEY (usuario_id),
  KEY usuarios_roles_fk (rol_id),
  CONSTRAINT usuarios_roles_fk
    FOREIGN KEY (rol_id)
    REFERENCES roles (role_id)
);

CREATE TABLE movimientos (
  movimiento_id INT NOT NULL AUTO_INCREMENT,
  usuario_id INT NOT NULL,
  proveedor_id INT DEFAULT NULL,
  tipo ENUM('entrada','salida') NOT NULL,
  fecha DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (movimiento_id),
  KEY fk_movimientos_usuario (usuario_id),
  KEY fk_movimientos_proveedor (proveedor_id),
  CONSTRAINT fk_movimientos_usuario
    FOREIGN KEY (usuario_id)
    REFERENCES usuarios (usuario_id)
    ON DELETE RESTRICT,
  CONSTRAINT fk_movimientos_proveedor
    FOREIGN KEY (proveedor_id)
    REFERENCES proveedores (proveedor_id)
    ON DELETE SET NULL
);

CREATE TABLE movimiento_detalle (
  detalle_id INT NOT NULL AUTO_INCREMENT,
  movimiento_id INT NOT NULL,
  insumo_id INT NOT NULL,
  cantidad INT NOT NULL DEFAULT 0,
  PRIMARY KEY (detalle_id),
  KEY movimiento_detalle_insumos_fk (insumo_id),
  KEY movimiento_detalle_movimientos_fk (movimiento_id),
  CONSTRAINT movimiento_detalle_insumos_fk
    FOREIGN KEY (insumo_id)
    REFERENCES insumos (insumo_id),
  CONSTRAINT movimiento_detalle_movimientos_fk
    FOREIGN KEY (movimiento_id)
    REFERENCES movimientos (movimiento_id)
);
CREATE TABLE categories (
  category_id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  PRIMARY KEY (category_id)
);

CREATE TABLE clients (
  client_id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(300) NOT NULL,
  contact_phone INT NOT NULL,
  contact_email VARCHAR(300) NOT NULL,
  PRIMARY KEY (client_id)
);

CREATE TABLE suppliers (
  supplier_id INT NOT NULL AUTO_INCREMENT,
  company_name VARCHAR(300) NOT NULL,
  company_phone INT NOT NULL,
  company_email VARCHAR(200) NOT NULL,
  contact_name VARCHAR(100) DEFAULT NULL,
  contact_phone INT DEFAULT NULL,
  contact_email VARCHAR(300) DEFAULT NULL,
  PRIMARY KEY (supplier_id)
);

CREATE TABLE roles (
  role_id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(300) NOT NULL,
  PRIMARY KEY (role_id)
);

CREATE TABLE users (
  user_id INT NOT NULL AUTO_INCREMENT,
  role_id INT NOT NULL,
  name VARCHAR(300) NOT NULL,
  password VARCHAR(100) NOT NULL,
  email VARCHAR(150) DEFAULT NULL,
  PRIMARY KEY (user_id),
  KEY fk_users_roles (role_id),
  CONSTRAINT fk_users_roles
    FOREIGN KEY (role_id)
    REFERENCES roles (role_id)
);

CREATE TABLE supplies (
  supply_id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(300) NOT NULL,
  type ENUM('physical','gastronomic') NOT NULL,
  stock INT NOT NULL DEFAULT 0,
  category_id INT NOT NULL,
  price DECIMAL(10,2) NOT NULL,
  image VARCHAR(300) NOT NULL,
  PRIMARY KEY (supply_id),
  KEY fk_supplies_category (category_id),
  CONSTRAINT fk_supplies_category
    FOREIGN KEY (category_id)
    REFERENCES categories (category_id)
);

CREATE TABLE purchases (
  purchase_id INT NOT NULL AUTO_INCREMENT,
  user_id INT NOT NULL,
  supplier_id INT NOT NULL,
  date DATETIME DEFAULT CURRENT_TIMESTAMP,
  subtotal DECIMAL(10,2) NOT NULL,
  discount DECIMAL(10,2) DEFAULT NULL,
  IVA DECIMAL(10,2) NOT NULL,
  total DECIMAL(10,2) NOT NULL,
  PRIMARY KEY (purchase_id),
  KEY fk_purchases_user (user_id),
  KEY fk_purchases_supplier (supplier_id),
  CONSTRAINT fk_purchases_user
    FOREIGN KEY (user_id)
    REFERENCES users (user_id),
  CONSTRAINT fk_purchases_supplier
    FOREIGN KEY (supplier_id)
    REFERENCES suppliers (supplier_id)
);

CREATE TABLE purchase_details (
  detail_id INT NOT NULL AUTO_INCREMENT,
  purchase_id INT NOT NULL,
  supply_id INT NOT NULL,
  quantity INT NOT NULL DEFAULT 0,
  unit_price DECIMAL(10,2) NOT NULL,
  discount DECIMAL(10,2) DEFAULT NULL,
  subtotal DECIMAL(10,2) NOT NULL,
  PRIMARY KEY (detail_id),
  KEY fk_purchase_details_purchase (purchase_id),
  KEY fk_purchase_details_supply (supply_id),
  CONSTRAINT fk_purchase_details_purchase
    FOREIGN KEY (purchase_id)
    REFERENCES purchases (purchase_id)
    ON DELETE CASCADE,
  CONSTRAINT fk_purchase_details_supply
    FOREIGN KEY (supply_id)
    REFERENCES supplies (supply_id)
);

CREATE TABLE sales (
  sale_id INT NOT NULL AUTO_INCREMENT,
  user_id INT NOT NULL,
  client_id INT NOT NULL,
  date DATETIME DEFAULT CURRENT_TIMESTAMP,
  subtotal DECIMAL(10,2) NOT NULL,
  discount DECIMAL(10,2) DEFAULT NULL,
  IVA DECIMAL(10,2) NOT NULL,
  total DECIMAL(10,2) NOT NULL,
  PRIMARY KEY (sale_id),
  KEY fk_sales_user (user_id),
  KEY fk_sales_client (client_id),
  CONSTRAINT fk_sales_user
    FOREIGN KEY (user_id)
    REFERENCES users (user_id),
  CONSTRAINT fk_sales_client
    FOREIGN KEY (client_id)
    REFERENCES clients (client_id)
);

CREATE TABLE sale_details (
  detail_id INT NOT NULL AUTO_INCREMENT,
  sale_id INT NOT NULL,
  supply_id INT NOT NULL,
  quantity INT NOT NULL DEFAULT 0,
  unit_price DECIMAL(10,2) NOT NULL,
  discount DECIMAL(10,2) DEFAULT NULL,
  subtotal DECIMAL(10,2) NOT NULL,
  PRIMARY KEY (detail_id),
  KEY fk_sale_details_sale (sale_id),
  KEY fk_sale_details_supply (supply_id),
  CONSTRAINT fk_sale_details_sale
    FOREIGN KEY (sale_id)
    REFERENCES sales (sale_id)
    ON DELETE CASCADE,
  CONSTRAINT fk_sale_details_supply
    FOREIGN KEY (supply_id)
    REFERENCES supplies (supply_id)
);
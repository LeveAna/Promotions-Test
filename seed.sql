CREATE DATABASE IF NOT EXISTS promotions_db;

USE promotions_db;

CREATE TABLE IF NOT EXISTS products (
    sku VARCHAR(20) PRIMARY KEY,
    name VARCHAR(255),
    category VARCHAR(50),
    price INT
);

INSERT INTO products (sku, name, category, price) VALUES
('000001', 'BV Lean leather ankle boots', 'boots', 89000),
('000002', 'BV Lean leather ankle boots', 'boots', 99000),
('000003', 'Ashlington leather ankle boots', 'boots', 71000),
('000004', 'Naima embellished suede sandals', 'sandals', 79500),
('000005', 'Nathane leather sneakers', 'sneakers', 59000);

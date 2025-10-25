-- Test: INSERT and UPDATE statements
INSERT INTO products (name, price, category) VALUES ('Laptop', 999.99, 'Electronics');
INSERT INTO products (name, price, category) VALUES ('Mouse', 29.99, 'Electronics');
INSERT INTO products (name, price, category) VALUES ('Monitor', 139.99, 'Electronics');
UPDATE users SET last_login = NOW() WHERE id = 123;
UPDATE orders SET status = 'shipped' WHERE order_date < '2023-12-01';

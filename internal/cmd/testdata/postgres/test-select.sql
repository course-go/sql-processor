-- Test: SELECT statements
SELECT * FROM users WHERE active = true;
SELECT id, name, email FROM customers WHERE created_at > '2023-01-01';
SELECT COUNT(*) FROM orders WHERE status = 'completed';
SELECT COUNT(*) FROM orders WHERE status = 'reserved';

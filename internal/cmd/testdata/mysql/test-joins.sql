-- Test: Complex queries with JOINs
SELECT u.name, o.total_amount, o.order_date
FROM users u
JOIN orders o ON u.id = o.user_id
WHERE o.status = 'completed'
ORDER BY o.order_date DESC;

SELECT p.name, c.name as category_name, COUNT(oi.id) as order_count
FROM products p
JOIN categories c ON p.category_id = c.id
LEFT JOIN order_items oi ON p.id = oi.product_id
GROUP BY p.id, p.name, c.name
HAVING COUNT(oi.id) > 5;

-- This is a comment
SELECT * FROM users;

INSERT INTO users (name, email)
VALUES ('John', 'john@example.com');


-- Another comment
UPDATE users
SET name = 'Jane'
WHERE id = 1;


UPDATE users
    SET name = 'Bob'
    WHERE id = 4;


DELETE FROM users WHERE id = 2;

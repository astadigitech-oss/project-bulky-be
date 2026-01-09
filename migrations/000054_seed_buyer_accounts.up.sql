-- migrations/000054_seed_buyer_accounts.up.sql
-- Seed initial buyer accounts for testing

-- Note: Password is hashed using bcrypt (cost=10) for "buyer123"
-- Hash generated: $2a$10$1Lswlh3jccTAipAsQUqkpegPuQ1xHGY2LtJqtFqpi8Jgm9CEhfXsm

INSERT INTO buyer (id, nama, username, email, password, telepon, is_active, is_verified, created_at, updated_at) VALUES
(
    uuid_generate_v4(),
    'John Doe',
    'johndoe',
    'john.doe@example.com',
    '$2a$10$1Lswlh3jccTAipAsQUqkpegPuQ1xHGY2LtJqtFqpi8Jgm9CEhfXsm',
    '081234567890',
    true,
    false,
    NOW(),
    NOW()
),
(
    uuid_generate_v4(),
    'Jane Smith',
    'janesmith',
    'jane.smith@example.com',
    '$2a$10$1Lswlh3jccTAipAsQUqkpegPuQ1xHGY2LtJqtFqpi8Jgm9CEhfXsm',
    '081234567891',
    true,
    false,
    NOW(),
    NOW()
),
(
    uuid_generate_v4(),
    'Bob Wilson',
    'bobwilson',
    'bob.wilson@example.com',
    '$2a$10$1Lswlh3jccTAipAsQUqkpegPuQ1xHGY2LtJqtFqpi8Jgm9CEhfXsm',
    '081234567892',
    true,
    false,
    NOW(),
    NOW()
);

-- Info Login:
-- Username/Email: johndoe / john.doe@example.com - Password: buyer123
-- Username/Email: janesmith / jane.smith@example.com - Password: buyer123
-- Username/Email: bobwilson / bob.wilson@example.com - Password: buyer123

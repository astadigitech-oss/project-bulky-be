-- migrations/000054_seed_buyer_accounts.down.sql
-- Rollback: Remove seeded buyer accounts

DELETE FROM buyer WHERE email IN (
    'john.doe@example.com',
    'jane.smith@example.com',
    'bob.wilson@example.com'
);

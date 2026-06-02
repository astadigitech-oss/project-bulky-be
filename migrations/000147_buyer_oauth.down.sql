-- migrations/000147_buyer_oauth.down.sql

DROP INDEX IF EXISTS buyer_oauth_provider_uid_unique;
DROP INDEX IF EXISTS idx_buyer_oauth_buyer_id;
DROP TABLE IF EXISTS buyer_oauth;

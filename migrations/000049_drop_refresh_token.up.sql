-- Drop refresh_token table karena tidak dipakai lagi
-- (Session 24 jam dengan single token)
DROP TABLE IF EXISTS refresh_token;

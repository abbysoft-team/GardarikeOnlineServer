ALTER TABLE IF EXISTS accounts
ADD COLUMN is_online bool NOT NULL DEFAULT false,
ADD COLUMN last_session_id varchar(36) NOT NULL DEFAULT ''

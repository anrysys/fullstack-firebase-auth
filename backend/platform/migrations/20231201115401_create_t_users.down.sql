DROP TRIGGER update_users_updated_at ON users;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS _update_all_updated_at;
-- DROP TYPE IF EXISTS _users_user_status_enum;
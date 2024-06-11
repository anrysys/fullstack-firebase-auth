DO $$
BEGIN
    -- Checking if a trigger exists on the table
    IF EXISTS (
        SELECT 1
        FROM pg_trigger
        WHERE not tgisinternal  -- Eliminating internal triggers
        AND tgname = 'update_users_updated_at'  -- Trigger name
    ) THEN
        -- Remove the trigger if it exists
        EXECUTE 'DROP TRIGGER update_users_updated_at ON users;';
    END IF;
END
$$;

-- Drop indexes
DROP INDEX IF EXISTS users_user_status_idx;
DROP INDEX IF EXISTS users_user_role_idx;
DROP INDEX IF EXISTS users_created_idx;
DROP INDEX IF EXISTS users_updated_idx;

-- DROP TYPE IF EXISTS _users_user_status_enum;

-- Drop table
DROP TABLE IF EXISTS users CASCADE;
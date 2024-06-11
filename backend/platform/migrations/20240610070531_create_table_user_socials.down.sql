DO $$
BEGIN
    -- Checking if a trigger exists on the table
    IF EXISTS (
        SELECT 1
        FROM pg_trigger
        WHERE not tgisinternal  -- Eliminating internal triggers
        AND tgname = 'update_user_socials_updated_at'  -- Trigger name
    ) THEN
        -- Remove the trigger if it exists
        EXECUTE 'DROP TRIGGER update_user_socials_updated_at ON user_socials;';
    END IF;
END
$$;

DROP TABLE IF EXISTS user_socials CASCADE;
 
-- Delete indexes from the table user_socials_user_id_firebase_uid_provider_id_idx
DROP INDEX IF EXISTS user_socials_user_id_firebase_uid_provider_id_idx;
DROP INDEX IF EXISTS user_socials_user_id_idx;
DROP INDEX IF EXISTS user_socials_provider_id_idx;
DROP INDEX IF EXISTS user_socials_provider_user_id_idx; 
DROP INDEX IF EXISTS user_socials_firebase_uid_uidx;

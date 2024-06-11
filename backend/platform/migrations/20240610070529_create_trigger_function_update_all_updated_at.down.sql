DO $$
BEGIN
    -- Checking whether a function with a given name exists
    IF EXISTS (
        SELECT 1
        FROM pg_proc p
        JOIN pg_namespace n ON p.pronamespace = n.oid
        WHERE n.nspname = 'public'  -- Schema name
        AND p.proname = '_update_all_updated_at'  -- Function name
    ) THEN
        -- Remove the function if it exists
        EXECUTE 'DROP FUNCTION public._update_all_updated_at();';
    END IF;
END
$$;
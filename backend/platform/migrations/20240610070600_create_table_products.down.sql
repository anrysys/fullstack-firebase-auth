-- Drop triggers first
DO $$
BEGIN
    -- Checking if a trigger exists on the table
    IF EXISTS (
        SELECT 1
        FROM pg_trigger
        WHERE not tgisinternal  -- Eliminating internal triggers
        AND tgname = 'update_products_updated_at'  -- Trigger name
    ) THEN
        -- Remove the trigger if it exists
        EXECUTE 'DROP TRIGGER update_products_updated_at ON products;';
    END IF;
END
$$;

-- Then drop the table
DROP TABLE IF EXISTS products CASCADE;

-- Now it's safe to drop the type as its dependencies have been removed
DROP TYPE IF EXISTS _products_product_status_enum;
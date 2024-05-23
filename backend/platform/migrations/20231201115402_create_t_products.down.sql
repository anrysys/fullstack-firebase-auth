-- Delete tables
DROP TRIGGER update_products_updated_at ON products;
DROP TABLE IF EXISTS products;
DROP TYPE IF EXISTS _products_product_status_enum;
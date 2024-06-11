-- Create ENUM type for product_status
DROP TYPE IF EXISTS _products_product_status_enum;
CREATE TYPE _products_product_status_enum AS ENUM ('pending', 'active', 'suspended', 'blocked', 'rejected', 'deleted');
-- Create products table
CREATE TABLE products (
    product_id SERIAL8 NOT NULL,
    title varchar (255) NOT NULL,
    description TEXT NOT NULL,
    user_id INT8 NOT NULL,
    product_status _products_product_status_enum NOT NULL DEFAULT 'pending'::_products_product_status_enum,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	CONSTRAINT haks_pkey PRIMARY KEY (product_id),
	CONSTRAINT products_product_status_check CHECK ((product_status = ANY (ARRAY[
            'pending'::_products_product_status_enum, 
            'active'::_products_product_status_enum, 
            'suspended'::_products_product_status_enum, 
            'blocked'::_products_product_status_enum, 
            'rejected'::_products_product_status_enum, 
            'deleted'::_products_product_status_enum
        ]))),
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);

COMMENT ON TABLE products IS 'products table';

-- Add indexes
CREATE INDEX products_product_status_idx ON products (product_status);

-- Add index on created_at and updated_at columns
CREATE INDEX products_created_at_idx ON users USING btree (created_at);
CREATE INDEX products_updated_at_idx ON users USING btree (updated_at);

-- Create triggers
CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION _update_all_updated_at();
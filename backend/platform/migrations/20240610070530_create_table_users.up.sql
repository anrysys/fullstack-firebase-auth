DROP TYPE IF EXISTS _users_user_role_enum;
CREATE TYPE _users_user_role_enum AS ENUM ('admin', 'customer', 'supplier', 'courier');

DROP TYPE IF EXISTS _users_user_status_enum;
CREATE TYPE _users_user_status_enum AS ENUM ('pending', 'active', 'suspended', 'blocked', 'rejected', 'deleted');

-- Create users table
CREATE TABLE users (
    user_id SERIAL8 NOT NULL,
    user_social_id INT8 DEFAULT NULL,
    firebase_uid VARCHAR (100) DEFAULT NULL,
    provider_id VARCHAR (100) DEFAULT NULL,
    email VARCHAR (100) NOT NULL UNIQUE, 
    short_name VARCHAR (36) UNIQUE,   
    first_name VARCHAR (50),
    last_name VARCHAR (50),
    lang CHAR(2) DEFAULT 'en', 
    phone_number VARCHAR(20), 
    password VARCHAR (100) NOT NULL,    
    user_status _users_user_status_enum NOT NULL DEFAULT 'pending'::_users_user_status_enum,
    user_role _users_user_role_enum NOT NULL DEFAULT 'customer'::_users_user_role_enum,
    photo VARCHAR (100) NOT NULL DEFAULT 'default.png', 
    activated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,    
    verified_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    CONSTRAINT users_email_uidx UNIQUE (email),
    CONSTRAINT users_short_name_uidx UNIQUE (short_name),
	CONSTRAINT users_pkey PRIMARY KEY (user_id),
	CONSTRAINT users_user_status_check CHECK ((user_status = ANY (
        ARRAY['pending'::_users_user_status_enum, 
            'active'::_users_user_status_enum, 
            'suspended'::_users_user_status_enum, 
            'blocked'::_users_user_status_enum, 
            'rejected'::_users_user_status_enum, 
            'deleted'::_users_user_status_enum])
        )),
	CONSTRAINT users_user_role_check CHECK ((user_role = ANY (
        ARRAY['admin'::_users_user_role_enum, 
            'customer'::_users_user_role_enum, 
            'supplier'::_users_user_role_enum,
            'courier'::_users_user_role_enum])
        ))

);       

COMMENT ON TABLE users IS 'Users table';
COMMENT ON COLUMN users.activated_at IS 'Date when the customer activated the account';
COMMENT ON COLUMN users.lang IS 'Two letter country code. See standart https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes , https://www.loc.gov/standards/iso639-2/';
COMMENT ON COLUMN users.email IS 'Email address';
COMMENT ON COLUMN users.phone_number IS 'Phone number';
COMMENT ON COLUMN users.photo IS 'User photo';
COMMENT ON COLUMN users.user_status IS 'User status';
COMMENT ON COLUMN users.user_role IS 'User role';

-- Add indexes

-- Create complex index: firebase_uid, provider_id
CREATE INDEX users_firebase_uid_provider_id_idx ON users (firebase_uid, provider_id);
CREATE INDEX users_email_idx ON users (email);
CREATE INDEX users_user_status_idx ON users (user_status);
CREATE INDEX users_user_role_idx ON users (user_role);
CREATE INDEX users_created_idx ON users USING btree (created_at);
CREATE INDEX users_updated_idx ON users USING btree (updated_at);


-- Create triggers
CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION _update_all_updated_at();

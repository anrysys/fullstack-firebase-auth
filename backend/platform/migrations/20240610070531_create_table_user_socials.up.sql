CREATE TABLE user_socials (
    user_social_id SERIAL8 NOT NULL,
    user_id INT8 NOT NULL,
    firebase_uid VARCHAR (100) NOT NULL,
    provider_id VARCHAR (100) DEFAULT NULL,
    tenant_id VARCHAR (100) DEFAULT NULL,
    email VARCHAR (100) DEFAULT NULL,    
    email_verified BOOLEAN DEFAULT FALSE,
    phone_number VARCHAR (20) DEFAULT NULL,        
    is_anonymous BOOLEAN DEFAULT FALSE,
    username VARCHAR (36) DEFAULT NULL,
    photo_url VARCHAR (150) DEFAULT NULL,
    provider_data json DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    CONSTRAINT user_socials_pkey PRIMARY KEY (user_social_id),
    CONSTRAINT user_socials_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE
);

COMMENT ON TABLE user_socials IS 'Profile user data from social media providers (Google, Facebook, etc.) ';
COMMENT ON COLUMN user_socials.email_verified IS 'Email verified status';
COMMENT ON COLUMN user_socials.firebase_uid IS 'Firebase user unique id from firebase';

-- Add indexes
CREATE INDEX users_firebase_uid_uidx ON user_socials (firebase_uid);

-- Create complex index with fields: user_id, firebase_uid and provider_id
CREATE INDEX user_socials_user_id_firebase_uid_provider_id_idx ON user_socials (user_id, firebase_uid, provider_id);

-- Create triggers
CREATE TRIGGER update_user_socials_updated_at
BEFORE UPDATE ON user_socials
FOR EACH ROW
EXECUTE FUNCTION _update_all_updated_at();










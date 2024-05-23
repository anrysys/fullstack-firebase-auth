-- Create auth_login_attempts table
CREATE TABLE auth_login_attempts (
    auth_login_attempt_id SERIAL8 NOT NULL,
    email VARCHAR (100) NOT NULL,    
    otp_code CHAR (5) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
	CONSTRAINT auth_login_attempts_pkey PRIMARY KEY (auth_login_attempt_id)
);

COMMENT ON TABLE auth_login_attempts IS 'Кол-во попыток для введения OTP CODE юзером. Для одного email может быть добавлено не более 3-ех otp_code за последние 15 мин.';

-- Add index on created_at and updated_at columns
CREATE INDEX idx_auth_login_attempts_created_at ON users USING btree (created_at);
CREATE INDEX idx_auth_login_attempts_updated_at ON users USING btree (updated_at);

CREATE INDEX auth_login_attempts_email_idx ON users (email);
CREATE INDEX auth_login_attempts_email_otp_code_cidx ON users (email, otp_code);


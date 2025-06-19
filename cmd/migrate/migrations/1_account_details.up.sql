CREATE TYPE account_status AS ENUM ('active', 'inactive', 'banned', 'suspended', 'pending_verification', 'deleted');

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    account_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    salt VARCHAR(64) NOT NULL,
    
    status account_status NOT NULL DEFAULT 'pending_verification',
    ban_reason TEXT,
    failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    last_failed_login TIMESTAMP WITH TIME ZONE,
    locked_until TIMESTAMP WITH TIME ZONE,
    
    is_online BOOLEAN DEFAULT false,
    last_login TIMESTAMP WITH TIME ZONE,
    last_logout TIMESTAMP WITH TIME ZONE,
    
    account_level INTEGER DEFAULT 1,
    experience_points BIGINT DEFAULT 0,
    currency_gold BIGINT DEFAULT 0,
    currency_premium BIGINT DEFAULT 0,
    
    email_verification_token VARCHAR(255),
    email_verified_at TIMESTAMP WITH TIME ZONE,
    password_reset_token VARCHAR(255),
    password_reset_expires TIMESTAMP WITH TIME ZONE,
    two_factor_enabled BOOLEAN DEFAULT false,
    two_factor_secret VARCHAR(32),
    
    notification_email BOOLEAN DEFAULT true,
    notification_push BOOLEAN DEFAULT true,
    notification_in_game BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT valid_username CHECK (username ~ '^[a-zA-Z0-9_-]{3,21}$'),
    CONSTRAINT valid_email CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    
    display_name VARCHAR(100),
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    date_of_birth DATE,
    country VARCHAR(3),
    timezone VARCHAR(50),
    language VARCHAR(10) DEFAULT 'en',
    
    bio TEXT,
    avatar_url VARCHAR(500),
    website VARCHAR(255),
    social_links JSONB,
    
    is_public BOOLEAN DEFAULT true,
    show_email BOOLEAN DEFAULT false,
    show_birthday BOOLEAN DEFAULT false,
    show_location BOOLEAN DEFAULT false,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT valid_website CHECK (website IS NULL OR website ~ '^https?://.+'),
    CONSTRAINT valid_avatar_url CHECK (avatar_url IS NULL OR avatar_url ~ '^https?://.+'),
    CONSTRAINT valid_date_of_birth CHECK (date_of_birth IS NULL OR date_of_birth <= CURRENT_DATE - INTERVAL '13 years')
);

CREATE INDEX idx_accounts_username ON accounts(username);
CREATE INDEX idx_accounts_email ON accounts(email);
CREATE INDEX idx_accounts_account_id ON accounts(account_id);
CREATE INDEX idx_accounts_status ON accounts(status);
CREATE INDEX idx_accounts_created_at ON accounts(created_at);
CREATE INDEX idx_accounts_last_login ON accounts(last_login);

CREATE INDEX idx_profiles_account_id ON profiles(account_id);
CREATE INDEX idx_profiles_display_name ON profiles(display_name);
CREATE INDEX idx_profiles_country ON profiles(country);
CREATE INDEX idx_profiles_is_public ON profiles(is_public);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_accounts_updated_at 
    BEFORE UPDATE ON accounts 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_profiles_updated_at 
    BEFORE UPDATE ON profiles 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE FUNCTION generate_account_id()
RETURNS UUID AS $$
BEGIN
    RETURN gen_random_uuid();
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION create_profile_for_account()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO profiles (account_id) VALUES (NEW.id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER create_profile_trigger
    AFTER INSERT ON accounts
    FOR EACH ROW
    EXECUTE FUNCTION create_profile_for_account();

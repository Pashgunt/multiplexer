-- +goose Up
CREATE TABLE IF NOT EXISTS target_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    base_url VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT target_services_name_check CHECK (service_name ~ '^[a-z0-9_-]+$')
);

CREATE TABLE IF NOT EXISTS service_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID NOT NULL REFERENCES target_services(id) ON DELETE CASCADE,
    environment VARCHAR(50) NOT NULL,
    auth_type VARCHAR(50) NOT NULL,
    client_id VARCHAR(255),
    client_secret_encrypted TEXT,
    token_url VARCHAR(255),
    scopes VARCHAR(512),
    api_key_encrypted TEXT,
    api_key_header VARCHAR(255) DEFAULT 'X-API-Key',
    username_encrypted VARCHAR(255),
    password_encrypted VARCHAR(255),
    jwt_secret_encrypted TEXT,
    jwt_algorithm VARCHAR(20) DEFAULT 'HS256',
    custom_headers JSONB,
    timeout_ms INTEGER DEFAULT 5000,
    retry_count INTEGER DEFAULT 3,
    retry_delay_ms INTEGER DEFAULT 1000,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_used_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT unique_service_environment UNIQUE (service_id, environment),
    CONSTRAINT valid_auth_type CHECK (auth_type IN ('oauth2_client_credentials', 'api_key', 'basic_auth', 'jwt', 'none')),
    CONSTRAINT valid_environment CHECK (environment IN ('dev', 'staging', 'prod', 'test'))
);

CREATE TABLE IF NOT EXISTS service_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID NOT NULL REFERENCES target_services(id) ON DELETE CASCADE,
    environment VARCHAR(50) NOT NULL,
    access_token_encrypted TEXT NOT NULL,
    token_type VARCHAR(50) DEFAULT 'Bearer',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    issued_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    refresh_token_encrypted TEXT,
    token_hash VARCHAR(255),
    expires_in_seconds INTEGER,
    created_by_instance VARCHAR(255),
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_active_token UNIQUE (service_id, environment),
    CONSTRAINT expires_at_future CHECK (expires_at > issued_at)
);

CREATE TABLE IF NOT EXISTS token_refresh_locks (
    service_id UUID NOT NULL REFERENCES target_services(id) ON DELETE CASCADE,
    environment VARCHAR(50) NOT NULL,
    locked_by VARCHAR(255) NOT NULL,
    locked_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    version BIGINT DEFAULT 0,
    PRIMARY KEY (service_id, environment),
    CONSTRAINT lock_expires_future CHECK (expires_at > locked_at)
);

CREATE TABLE IF NOT EXISTS token_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID NOT NULL REFERENCES target_services(id),
    environment VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    old_token_hash VARCHAR(255),
    new_token_hash VARCHAR(255),
    old_expires_at TIMESTAMP WITH TIME ZONE,
    new_expires_at TIMESTAMP WITH TIME ZONE,
    reason VARCHAR(255),
    instance_id VARCHAR(255),
    latency_ms INTEGER,
    error_message TEXT,
    error_type VARCHAR(100),
    request_id VARCHAR(255),
    trace_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS token_metrics_aggr (
    service_id UUID NOT NULL REFERENCES target_services(id),
    environment VARCHAR(50) NOT NULL,
    date DATE NOT NULL,
    hour INTEGER NOT NULL,
    requests_total BIGINT DEFAULT 0,
    cache_hits BIGINT DEFAULT 0,
    cache_misses BIGINT DEFAULT 0,
    refreshes_total BIGINT DEFAULT 0,
    refreshes_success BIGINT DEFAULT 0,
    refreshes_failed BIGINT DEFAULT 0,
    refreshes_skipped BIGINT DEFAULT 0,
    avg_latency_ms INTEGER DEFAULT 0,
    p95_latency_ms INTEGER DEFAULT 0,
    p99_latency_ms INTEGER DEFAULT 0,
    max_latency_ms INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (service_id, environment, date, hour)
);

CREATE TABLE IF NOT EXISTS service_credentials_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    credential_id UUID NOT NULL,
    service_id UUID NOT NULL,
    environment VARCHAR(50) NOT NULL,
    auth_type VARCHAR(50) NOT NULL,
    client_id VARCHAR(255),
    client_secret_encrypted TEXT,
    token_url VARCHAR(255),
    changed_by VARCHAR(255),
    change_reason VARCHAR(255),
    change_type VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS token_health_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID NOT NULL REFERENCES target_services(id),
    environment VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    last_check_time TIMESTAMP WITH TIME ZONE NOT NULL,
    next_check_time TIMESTAMP WITH TIME ZONE NOT NULL,
    consecutive_failures INTEGER DEFAULT 0,
    last_error_message TEXT,
    expires_at TIMESTAMP WITH TIME ZONE,
    expires_in_seconds INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_service_health UNIQUE (service_id, environment)
);

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS token_health_checks;
DROP TABLE IF EXISTS service_credentials_history;
DROP TABLE IF EXISTS token_metrics_aggr;
DROP TABLE IF EXISTS token_audit_log;
DROP TABLE IF EXISTS token_refresh_locks;
DROP TABLE IF EXISTS service_tokens;
DROP TABLE IF EXISTS service_credentials;
DROP TABLE IF EXISTS target_services;
-- +goose StatementEnd
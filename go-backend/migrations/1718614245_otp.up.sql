CREATE TABLE otp (
    otp CHAR(6) PRIMARY KEY CHECK (otp ~ '^[0-9]{6}$'),
    org_id BIGINT REFERENCES organizations(id),
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC')
);
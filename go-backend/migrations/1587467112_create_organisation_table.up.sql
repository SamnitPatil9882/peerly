CREATE TABLE IF NOT EXISTS organizations (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name varchar(50),
    contact_email varchar(50) ,
    domain_name varchar(45) ,
    subscription_status INTEGER,
    subscription_valid_upto TIMESTAMP,
    hi5_limit INTEGER,
    hi5_quota_renewal_frequency VARCHAR(50),
    timezone VARCHAR(100) DEFAULT 'UTC',
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC'),
    created_by BIGINT REFERENCES users(id),
    soft_delete BOOLEAN DEFAULT FALSE,
    soft_delete_by BIGINT REFERENCES users(id),
    is_email_verified BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'UTC')
);

-- TODO Add some indexes

-- Make users(org_id) a foreign key to organizations(id)
ALTER TABLE IF EXISTS users ADD CONSTRAINT fk_users_org_id_organizations_id FOREIGN KEY (org_id) REFERENCES organizations(id);

-- Domain name should be a unique value
CREATE UNIQUE INDEX IF NOT EXISTS fk_organizations_domain_name_unique ON organizations(domain_name);

-- Create the Josh Software organization
-- TODO: Get proper values (management decisions) for these things...
INSERT INTO organizations (
    id, name, contact_email, domain_name, subscription_status, subscription_valid_upto, hi5_limit, 
    hi5_quota_renewal_frequency, timezone, created_at,created_by,soft_delete,soft_delete_by,updated_at
) VALUES (
    DEFAULT, 'Josh Software', 'j.austin.hughey@joshsoftware.com', 'joshsoftware.com', 9999, '2024-06-30 23:59:59', 9999, 
    '1 Week', DEFAULT, DEFAULT,0,false,0,DEFAULT
) ON CONFLICT DO NOTHING;

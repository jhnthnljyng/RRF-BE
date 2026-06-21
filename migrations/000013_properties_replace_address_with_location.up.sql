ALTER TABLE properties
    DROP COLUMN IF EXISTS address,
    DROP COLUMN IF EXISTS city,
    DROP COLUMN IF EXISTS state,
    DROP COLUMN IF EXISTS zip_code,
    ADD COLUMN location VARCHAR(255) NOT NULL DEFAULT '';

ALTER TABLE properties ALTER COLUMN location DROP DEFAULT;

DROP INDEX IF EXISTS idx_properties_city;

CREATE INDEX idx_properties_location ON properties(location);

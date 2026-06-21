DROP INDEX IF EXISTS idx_properties_location;

ALTER TABLE properties
    DROP COLUMN IF EXISTS location,
    ADD COLUMN address  TEXT NOT NULL DEFAULT '',
    ADD COLUMN city     VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN state    VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN zip_code VARCHAR(20);

ALTER TABLE properties
    ALTER COLUMN address DROP DEFAULT,
    ALTER COLUMN city    DROP DEFAULT,
    ALTER COLUMN state   DROP DEFAULT;

CREATE INDEX idx_properties_city ON properties(city);

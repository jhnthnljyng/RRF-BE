CREATE TYPE room_type AS ENUM ('single', 'shared', 'studio', 'apartment', 'house');
CREATE TYPE furnishing_type AS ENUM ('furnished', 'unfurnished', 'partial');

CREATE TABLE properties (
    id              SERIAL PRIMARY KEY,
    owner_id        INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    address         TEXT NOT NULL,
    city            VARCHAR(100) NOT NULL,
    state           VARCHAR(100) NOT NULL,
    zip_code        VARCHAR(20),
    latitude        DECIMAL(10, 8),
    longitude       DECIMAL(11, 8),
    monthly_rent    DECIMAL(10, 2) NOT NULL,
    room_type       room_type NOT NULL,
    furnishing      furnishing_type NOT NULL DEFAULT 'unfurnished',
    max_occupants   SMALLINT NOT NULL DEFAULT 1,
    available_from  DATE NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_properties_owner_id     ON properties(owner_id);
CREATE INDEX idx_properties_city         ON properties(city);
CREATE INDEX idx_properties_monthly_rent ON properties(monthly_rent);
CREATE INDEX idx_properties_available    ON properties(available_from);
CREATE INDEX idx_properties_is_active    ON properties(is_active);

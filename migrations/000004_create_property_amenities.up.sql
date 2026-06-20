CREATE TABLE property_amenities (
    id          SERIAL PRIMARY KEY,
    property_id INTEGER NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    UNIQUE(property_id, name)
);

CREATE INDEX idx_property_amenities_property_id ON property_amenities(property_id);

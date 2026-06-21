ALTER TABLE properties
    ADD COLUMN bedrooms  SMALLINT,
    ADD COLUMN bathrooms SMALLINT;

ALTER TYPE room_type ADD VALUE IF NOT EXISTS 'whole_unit';
ALTER TYPE room_type ADD VALUE IF NOT EXISTS 'room';

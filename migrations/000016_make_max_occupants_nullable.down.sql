UPDATE properties SET max_occupants = 1 WHERE max_occupants IS NULL;
ALTER TABLE properties ALTER COLUMN max_occupants SET NOT NULL;
ALTER TABLE properties ALTER COLUMN max_occupants SET DEFAULT 1;

DROP TABLE IF EXISTS roommate_profiles;
DROP TYPE IF EXISTS lifestyle_type;
DROP TYPE IF EXISTS cleanliness_type;

CREATE TYPE gender_preference_type AS ENUM ('male', 'female', 'any');

ALTER TABLE properties
    ADD COLUMN gender_preference gender_preference_type NOT NULL DEFAULT 'any';

ALTER TABLE properties DROP COLUMN IF EXISTS gender_preference;

DROP TYPE IF EXISTS gender_preference_type;

CREATE TYPE lifestyle_type AS ENUM ('early_bird', 'night_owl', 'flexible');
CREATE TYPE cleanliness_type AS ENUM ('very_clean', 'clean', 'moderate', 'relaxed');

CREATE TABLE roommate_profiles (
    id                  SERIAL PRIMARY KEY,
    user_id             INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    bio                 TEXT,
    occupation          VARCHAR(100),
    budget_min          DECIMAL(10, 2),
    budget_max          DECIMAL(10, 2),
    preferred_move_in   DATE,
    preferred_city      VARCHAR(100),
    lifestyle           lifestyle_type,
    cleanliness         cleanliness_type,
    smoking_ok          BOOLEAN NOT NULL DEFAULT FALSE,
    pets_ok             BOOLEAN NOT NULL DEFAULT FALSE,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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

CREATE INDEX idx_roommate_profiles_user_id        ON roommate_profiles(user_id);
CREATE INDEX idx_roommate_profiles_preferred_city ON roommate_profiles(preferred_city);
CREATE INDEX idx_roommate_profiles_budget         ON roommate_profiles(budget_min, budget_max);

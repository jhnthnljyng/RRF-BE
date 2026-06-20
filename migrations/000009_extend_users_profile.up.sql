CREATE TYPE cooking_frequency_type AS ENUM ('never', 'rarely', 'sometimes', 'often', 'always');

ALTER TABLE users
    ADD COLUMN smoking            SMALLINT NOT NULL DEFAULT 0 CHECK (smoking IN (0, 1)),
    ADD COLUMN pet_owner          SMALLINT NOT NULL DEFAULT 0 CHECK (pet_owner IN (0, 1)),
    ADD COLUMN pet_friendly       SMALLINT NOT NULL DEFAULT 0 CHECK (pet_friendly IN (0, 1)),
    ADD COLUMN bio                TEXT,
    ADD COLUMN gender             VARCHAR(50),
    ADD COLUMN occupation         VARCHAR(100),
    ADD COLUMN nationality        VARCHAR(100),
    ADD COLUMN cooking_frequency  cooking_frequency_type;

CREATE TABLE user_socials (
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    platform    VARCHAR(100) NOT NULL,
    url         TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, platform)
);

CREATE INDEX idx_user_socials_user_id ON user_socials(user_id);

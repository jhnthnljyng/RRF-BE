DROP TABLE IF EXISTS user_socials;

ALTER TABLE users
    DROP COLUMN IF EXISTS username,
    DROP COLUMN IF EXISTS smoking,
    DROP COLUMN IF EXISTS pet_owner,
    DROP COLUMN IF EXISTS pet_friendly,
    DROP COLUMN IF EXISTS bio,
    DROP COLUMN IF EXISTS gender,
    DROP COLUMN IF EXISTS occupation,
    DROP COLUMN IF EXISTS nationality,
    DROP COLUMN IF EXISTS cooking_frequency;

DROP TYPE IF EXISTS cooking_frequency_type;

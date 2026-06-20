ALTER TABLE users ALTER COLUMN role DROP DEFAULT;
ALTER TABLE users ALTER COLUMN role TYPE TEXT;

DROP TYPE user_role;
CREATE TYPE user_role AS ENUM ('seeker', 'landlord', 'both');

UPDATE users SET role = 'seeker'   WHERE role = 'tenant';
UPDATE users SET role = 'landlord' WHERE role = 'owner';

ALTER TABLE users ALTER COLUMN role TYPE user_role USING role::user_role;
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'seeker';

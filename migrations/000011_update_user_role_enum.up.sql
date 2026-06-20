ALTER TABLE users ALTER COLUMN role DROP DEFAULT;
ALTER TABLE users ALTER COLUMN role TYPE TEXT;

DROP TYPE user_role;
CREATE TYPE user_role AS ENUM ('owner', 'tenant');

UPDATE users SET role = 'tenant' WHERE role IN ('seeker', 'both');
UPDATE users SET role = 'owner'  WHERE role = 'landlord';

ALTER TABLE users ALTER COLUMN role TYPE user_role USING role::user_role;
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'tenant';

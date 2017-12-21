UPDATE users SET login = '' WHERE login IS NULL;
UPDATE users SET password_hash = '' WHERE password_hash IS NULL;
UPDATE users SET salt = '' WHERE salt IS NULL;
ALTER TABLE users
  ALTER COLUMN login SET NOT NULL,
  ALTER COLUMN password_hash SET NOT NULL,
  ALTER COLUMN salt SET NOT NULL;

DELETE FROM accounts WHERE user_id IS NULL; -- accounts must be linked to user
UPDATE accounts SET google = '' WHERE google IS NULL;
UPDATE accounts SET facebook = '' WHERE facebook IS NULL;
UPDATE accounts SET github = '' WHERE github IS NULL;
ALTER TABLE accounts
  ALTER COLUMN user_id SET NOT NULL,
  ALTER COLUMN google SET NOT NULL,
  ALTER COLUMN facebook SET NOT NULL,
  ALTER COLUMN github SET NOT NULL;

DELETE FROM links WHERE user_id IS NULL or type IS NULL; -- same as before
UPDATE links SET is_active = FALSE WHERE is_active IS NULL;
ALTER TABLE links
  ALTER COLUMN user_id SET NOT NULL,
  ALTER COLUMN type SET NOT NULL,
  ALTER COLUMN is_active SET NOT NULL;

DELETE FROM profiles WHERE user_id IS NULL;
UPDATE profiles SET referral = '' WHERE referral IS NULL;
UPDATE profiles SET access = '' WHERE access IS NULL;
UPDATE profiles SET data = '{}' WHERE data IS NULL;
ALTER TABLE profiles
  ALTER COLUMN user_id SET NOT NULL,
  ALTER COLUMN referral SET NOT NULL,
  ALTER COLUMN access SET NOT NULL,
  ALTER COLUMN data SET NOT NULL,
  ALTER COLUMN data SET DEFAULT '{}';

DELETE FROM tokens WHERE user_id IS NULL OR session_id IS NULL;
UPDATE tokens SET is_active = FALSE WHERE is_active IS NULL;
ALTER TABLE tokens
  ALTER COLUMN user_id SET NOT NULL,
  ALTER COLUMN is_active SET NOT NULL,
  ALTER COLUMN user_id SET NOT NULL;
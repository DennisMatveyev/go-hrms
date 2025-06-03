CREATE TABLE IF NOT EXISTS profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    first_name VARCHAR(32) NOT NULL,
	last_name VARCHAR(32) NOt NULL,
	phone_number VARCHAR(20) NOT NULL,
	photo_path VARCHAR(20),
	header VARCHAR(200)
);
CREATE INDEX idx_profiles_user_id ON profiles(user_id);
ALTER TABLE profiles ADD CONSTRAINT fk_profiles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    company_name VARCHAR(50) NOT NULL,
	position VARCHAR(50) NOT NULL,
	start_date DATE NOT NULL,
    end_date DATE,
    description TEXT NOT NULL
);
CREATE INDEX idx_jobs_user_id ON jobs(user_id);
ALTER TABLE jobs ADD CONSTRAINT fk_jobs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

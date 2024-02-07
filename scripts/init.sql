-- Creating tables
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE,
    first_name TEXT,
    last_name TEXT
);

CREATE TABLE IF NOT EXISTS candidates (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE,
    current_position TEXT,
    resume TEXT,
    bio TEXT
);

CREATE TABLE IF NOT EXISTS recruiters (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE,
    company_public_id UUID,
    int INT
);

CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE,
    recruiter_id INT,
    name TEXT
);

CREATE TABLE IF NOT EXISTS positions (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE,
    name TEXT,
    recruiters_public_id UUID UNIQUE
);

CREATE TABLE IF NOT EXISTS skills (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE,
    name TEXT
);

CREATE TABLE IF NOT EXISTS areas (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE,
    position_id INT,
    name TEXT
);

CREATE TABLE IF NOT EXISTS interviews (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE,
    score INT
);

CREATE TABLE IF NOT EXISTS videos (
    id SERIAL PRIMARY KEY,
    public_id UUID UNIQUE,
    interviews_public_id UUID UNIQUE,
    path TEXT
);

CREATE TABLE IF NOT EXISTS auth (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE,
    login TEXT UNIQUE,
    password TEXT
);

CREATE TABLE IF NOT EXISTS positions_skills (
    position_id INT,
    skills_id INT,
    PRIMARY KEY (position_id, skills_id)
);

CREATE TABLE IF NOT EXISTS candidate_skills (
    candidate_id INT,
    skills_id INT,
    PRIMARY KEY (candidate_id, skills_id)
);

CREATE TABLE IF NOT EXISTS user_interviews (
    candidate_id INT,
    recruiter_id INT,
    interview_id INT,
    PRIMARY KEY (candidate_id, recruiter_id, interview_id)
);

-- Creating references
ALTER TABLE users ADD CONSTRAINT fk_users_candidates FOREIGN KEY (public_id) REFERENCES candidates(public_id);
ALTER TABLE users ADD CONSTRAINT fk_users_recruiters FOREIGN KEY (public_id) REFERENCES recruiters(public_id);
ALTER TABLE user_interviews ADD CONSTRAINT fk_user_interviews_candidates FOREIGN KEY (candidate_id) REFERENCES candidates(id);
ALTER TABLE user_interviews ADD CONSTRAINT fk_user_interviews_recruiters FOREIGN KEY (recruiter_id) REFERENCES recruiters(id);
ALTER TABLE user_interviews ADD CONSTRAINT fk_user_interviews_interviews FOREIGN KEY (interview_id) REFERENCES interviews(id);
ALTER TABLE companies ADD CONSTRAINT fk_companies_recruiters FOREIGN KEY (recruiter_id) REFERENCES recruiters(id);
ALTER TABLE auth ADD CONSTRAINT fk_auth_users FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE positions_skills ADD CONSTRAINT fk_positions_skills_skills FOREIGN KEY (skills_id) REFERENCES skills(id);
ALTER TABLE positions_skills ADD CONSTRAINT fk_positions_skills_positions FOREIGN KEY (position_id) REFERENCES positions(id);
ALTER TABLE candidate_skills ADD CONSTRAINT fk_candidate_skills_skills FOREIGN KEY (skills_id) REFERENCES skills(id);
ALTER TABLE candidate_skills ADD CONSTRAINT fk_candidate_skills_candidates FOREIGN KEY (candidate_id) REFERENCES candidates(id);
ALTER TABLE positions ADD CONSTRAINT fk_recruiters_positions FOREIGN KEY (recruiters_public_id) REFERENCES recruiters(public_id);
ALTER TABLE videos ADD CONSTRAINT fk_interviews_videos FOREIGN KEY (interviews_public_id) REFERENCES interviews(public_id);

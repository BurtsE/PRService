CREATE TABLE IF NOT EXISTS teams (
    name TEXT PRIMARY KEY
);

-- CREATE TABLE IF NOT EXISTS team_members (
--     team_name TEXT NOT NULL REFERENCES teams(name) ON DELETE CASCADE,
--     user_id VARCHAR(64) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--     PRIMARY KEY (team_name, user_id)
-- );
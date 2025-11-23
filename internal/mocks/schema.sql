CREATE TABLE IF NOT EXISTS teams
(
    name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users
(
    id        VARCHAR(255) PRIMARY KEY,
    name      VARCHAR(255) NOT NULL,
    team_name VARCHAR(255) REFERENCES teams (name),
    is_active BOOLEAN      NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS pull_requests
(
    id         VARCHAR(255) PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    author_id  VARCHAR(255) REFERENCES users (id),
    status     VARCHAR(50)  NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    merged_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pull_request_reviewers
(
    pull_request_id VARCHAR(255) REFERENCES pull_requests (id) ON DELETE CASCADE,
    user_id         VARCHAR(255) REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (pull_request_id, user_id)
);

-- Add index for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_team_name ON users (team_name);
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_user_id ON pull_request_reviewers (user_id);
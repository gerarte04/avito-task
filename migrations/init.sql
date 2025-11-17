CREATE TABLE teams (
    name        varchar(100)    PRIMARY KEY
);

CREATE TABLE users (
    id          varchar(100)    PRIMARY KEY,
    name        varchar(100)    NOT NULL,
    team_name   varchar(100)    REFERENCES teams(name) ON DELETE SET NULL ON UPDATE CASCADE,
    is_active   bool            NOT NULL
);

CREATE INDEX users_team_name_idx ON users(team_name, is_active);

CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');
CREATE TABLE pull_requests (
    id          varchar(100)    PRIMARY KEY,
    name        varchar(100)    NOT NULL,
    author_id   varchar(100)    REFERENCES users(id),
    status      pr_status       NOT NULL DEFAULT 'OPEN',
    created_at  timestamp       DEFAULT CURRENT_TIMESTAMP,
    merged_at   timestamp
);

CREATE INDEX prs_status_id_idx ON pull_requests(status, id);

CREATE TABLE reviewers (
    pr_id   varchar(100) NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    user_id varchar(100) NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX reviewers_pr_idx ON reviewers(pr_id);
CREATE INDEX reviewers_user_pr_idx ON reviewers(user_id, pr_id);

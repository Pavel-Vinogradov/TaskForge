CREATE TABLE team_members
(
    user_id   INT         NOT NULL,
    team_id   INT         NOT NULL,
    role      VARCHAR(20) NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, team_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (team_id) REFERENCES teams (id)
);

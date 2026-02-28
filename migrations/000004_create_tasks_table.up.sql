CREATE TABLE tasks
(
    id          INT PRIMARY KEY AUTO_INCREMENT,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    status      VARCHAR(20)  NOT NULL,
    assignee_id INT          NOT NULL,
    team_id     INT          NOT NULL,
    created_by  INT          NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (assignee_id) REFERENCES users (id),
    FOREIGN KEY (team_id) REFERENCES teams (id),
    FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE task_history
(
    id         INT PRIMARY KEY AUTO_INCREMENT,
    task_id    INT          NOT NULL,
    changed_by INT          NOT NULL,
    field_name VARCHAR(255) NOT NULL,
    old_value  TEXT,
    new_value  TEXT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks (id),
    FOREIGN KEY (changed_by) REFERENCES users (id)
);
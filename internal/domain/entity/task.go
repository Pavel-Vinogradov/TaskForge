package entity

import "time"

type Task struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	TeamID      int        `json:"team_id"`
	CreatedBy   int        `json:"created_by"`
	AssigneeID  int        `json:"assignee_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

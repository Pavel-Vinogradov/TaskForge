package entity

import "time"

type Task struct {
	ID          int
	Title       string
	Description string
	Status      TaskStatus
	TeamID      int
	CreatedBy   int
	AssigneeID  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

package entity

import "time"

type TaskHistory struct {
	ID        int
	TaskID    int
	Field     string
	OldValue  string
	NewValue  string
	ChangedBy int
	ChangedAt time.Time
}

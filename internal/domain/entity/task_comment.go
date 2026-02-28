package entity

import "time"

type TaskComment struct {
	ID        int
	TaskID    int
	UserID    int
	Comment   string
	CreatedAt time.Time
}

package entity

import "time"

type TaskHistory struct {
	Id        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	Field     string    `json:"field"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	ChangedBy int       `json:"changed_by"`
	ChangedAt time.Time `json:"changed_at"`
}

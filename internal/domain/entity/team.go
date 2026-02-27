package entity

import "time"

type Team struct {
	ID        int
	Name      string
	CreatedBy int
	CreatedAt time.Time
	UpdatedAt time.Time
}

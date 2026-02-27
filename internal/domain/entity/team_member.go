package entity

import "time"

type TeamMember struct {
	TeamID   int
	UserID   int
	Role     Role
	JoinedAt time.Time
}

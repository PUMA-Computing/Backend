package models

import "time"

type Roles struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

package types

import "time"

type UserParams struct {
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

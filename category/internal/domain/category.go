package domain

import "time"

type CategoryParams struct {
	Uuid      string
	Name      string
	Image     []byte
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

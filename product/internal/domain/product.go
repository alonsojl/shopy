package domain

import (
	"strconv"
	"time"
)

type ProductParams struct {
	Uuid      string
	Name      string
	Price     float64
	QRCode    string
	IsTop     bool
	Location  string
	Category  Category
	Image     []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Category struct {
	Uuid string
	Name string
}

func Top(s string) (int, error) {
	if s == "" {
		// default to 0
		return 0, nil
	}
	return strconv.Atoi(s)
}

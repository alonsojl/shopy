package models

type Products []*Product
type Product struct {
	Uuid      string   `json:"uuid"`
	Name      string   `json:"name"`
	Price     float64  `json:"price"`
	Image     string   `json:"image"`
	QRCode    string   `json:"qrcode"`
	IsTop     bool     `json:"is_top"`
	Category  Category `json:"category"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type Category struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

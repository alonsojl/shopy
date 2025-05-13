package dynamodb

type ProductTable struct {
	Uuid         string  `dynamodbav:"uuid"`
	Name         string  `dynamodbav:"name"`
	Price        float64 `dynamodbav:"price"`
	Image        string  `dynamodbav:"image"`
	QRCode       string  `dynamodbav:"qrcode"`
	IsTop        bool    `dynamodbav:"is_top"`
	CategoryUuid string  `dynamodbav:"category_uuid"`
	CategoryName string  `dynamodbav:"category_name"`
	CreatedAt    string  `dynamodbav:"created_at"`
	UpdatedAt    string  `dynamodbav:"updated_at"`
}

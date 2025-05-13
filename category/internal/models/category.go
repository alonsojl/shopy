package models

type Categories []*Category
type Category struct {
	Uuid      string `json:"uuid" dynamodbav:"uuid"`
	Name      string `json:"name" dynamodbav:"name"`
	Image     string `json:"image" dynamodbav:"image"`
	CreatedAt string `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt string `json:"updated_at" dynamodbav:"updated_at"`
}

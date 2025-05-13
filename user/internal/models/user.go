package models

type Users []*User
type User struct {
	Email     string `json:"email" dynamodbav:"email"`
	Password  string `json:"password" dynamodbav:"password"`
	CreatedAt string `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt string `json:"updated_at" dynamodbav:"updated_at"`
}

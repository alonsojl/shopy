package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"shopy/internal/models"
	"time"

	mtypes "shopy/internal/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const tableName = "user"

type User struct {
	logger *slog.Logger
	client *dynamodb.Client
}

func NewUser(logger *slog.Logger, client *dynamodb.Client) *User {
	return &User{
		logger: logger,
		client: client,
	}
}

func (u *User) AddUser(ctx context.Context, params mtypes.UserParams) (*models.User, error) {
	var (
		user = &models.User{
			Email:     params.Email,
			Password:  params.Password,
			CreatedAt: params.CreatedAt.Format(time.DateTime),
			UpdatedAt: params.UpdatedAt.Format(time.DateTime),
		}
	)

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, fmt.Errorf("error marshaling item: %w", err)
	}

	_, err = u.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("error adding item: %w", err)
	}

	return user, nil
}

func (u *User) DelUser(ctx context.Context, email string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"email": &types.AttributeValueMemberS{Value: email},
		},
		ConditionExpression: aws.String("attribute_exists(email)"),
	}

	_, err := u.client.DeleteItem(ctx, input)
	if err != nil {
		var errf *types.ConditionalCheckFailedException
		if errors.As(err, &errf) {
			return mtypes.ErrNotFound
		}

		return fmt.Errorf("error deleting item: %w", err)
	}

	return nil
}

func (u *User) GetUser(ctx context.Context, email string) (*models.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"email": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := u.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting item: %w", err)
	}

	var user models.User
	if err = attributevalue.UnmarshalMap(result.Item, &user); err != nil {
		return nil, fmt.Errorf("error unmarshaling item: %w", err)
	}

	return &user, nil
}

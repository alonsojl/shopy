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

type Category struct {
	logger *slog.Logger
	client *dynamodb.Client
	table  string
}

func NewCategory(logger *slog.Logger, client *dynamodb.Client) *Category {
	return &Category{
		logger: logger,
		client: client,
		table:  "category",
	}
}

func (c *Category) GetCategories(ctx context.Context) (models.Categories, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(c.table),
	}

	result, err := c.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	var categories models.Categories
	if err = attributevalue.UnmarshalListOfMaps(result.Items, &categories); err != nil {
		return nil, fmt.Errorf("error unmarshaling items: %w", err)
	}

	return categories, nil
}

func (c *Category) AddCategory(ctx context.Context, params mtypes.CategoryParams) (*models.Category, error) {
	category := &models.Category{
		Uuid:      params.Uuid,
		Name:      params.Name,
		Image:     params.Location,
		CreatedAt: params.CreatedAt.Format(time.DateTime),
		UpdatedAt: params.UpdatedAt.Format(time.DateTime),
	}

	item, err := attributevalue.MarshalMap(category)
	if err != nil {
		return nil, fmt.Errorf("error marshaling item: %w", err)
	}

	_, err = c.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(c.table),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("error adding item: %w", err)
	}

	return category, nil
}

func (c *Category) DelCategory(ctx context.Context, uuid string) (*models.Category, error) {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(c.table),
		Key: map[string]types.AttributeValue{
			"uuid": &types.AttributeValueMemberS{Value: uuid},
		},
		ConditionExpression: aws.String("attribute_exists(#uuid)"),
		ExpressionAttributeNames: map[string]string{
			"#uuid": "uuid",
		},
		ReturnValues: types.ReturnValueAllOld,
	}

	result, err := c.client.DeleteItem(ctx, input)
	if err != nil {
		var errf *types.ConditionalCheckFailedException
		if errors.As(err, &errf) {
			return nil, mtypes.ErrNotFound
		}

		return nil, fmt.Errorf("error deleting item: %w", err)
	}

	var category models.Category
	if err = attributevalue.UnmarshalMap(result.Attributes, &category); err != nil {
		return nil, fmt.Errorf("error marshaling item: %w", err)
	}

	return &category, nil
}

package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"shopy/internal/models"
	"strconv"
	"time"

	mtypes "shopy/internal/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Product struct {
	logger    *slog.Logger
	client    *dynamodb.Client
	tableName string
}

func NewProduct(logger *slog.Logger, client *dynamodb.Client) *Product {
	return &Product{
		logger:    logger,
		client:    client,
		tableName: "product",
	}
}

func (p *Product) GetProductsByCategory(ctx context.Context, uuid string) (models.Products, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(p.tableName),
		IndexName:              aws.String("GSI_CATEGORY"),
		KeyConditionExpression: aws.String("category_uuid = :uuid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uuid": &types.AttributeValueMemberS{Value: uuid},
		},
		ScanIndexForward: aws.Bool(true),
	}

	result, err := p.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	if len(result.Items) == 0 {
		return models.Products{}, nil
	}

	var product ProductTable
	products := make(models.Products, len(result.Items))
	for i, item := range result.Items {
		if err = attributevalue.UnmarshalMap(item, &product); err != nil {
			return nil, fmt.Errorf("error unmarshaling item: %w", err)
		}
		products[i] = assembleProduct(product)
	}

	return products, nil
}

func (p *Product) GetProductsByQRCode(ctx context.Context, qrcode string) (models.Products, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(p.tableName),
		IndexName:              aws.String("GSI_QRCODE"),
		KeyConditionExpression: aws.String("qrcode = :qrcode"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":qrcode": &types.AttributeValueMemberS{Value: qrcode},
		},
	}

	result, err := p.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	if len(result.Items) == 0 {
		return models.Products{}, nil
	}

	var product ProductTable
	products := make(models.Products, len(result.Items))

	for i, item := range result.Items {
		if err = attributevalue.UnmarshalMap(item, &product); err != nil {
			return nil, fmt.Errorf("error unmarshaling item: %w", err)
		}

		products[i] = assembleProduct(product)
	}

	return products, nil
}

func (p *Product) GetProductsByName(ctx context.Context, name string) (models.Products, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String(p.tableName),
		FilterExpression: aws.String("begins_with(#name, :name)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{Value: name},
		},
		ExpressionAttributeNames: map[string]string{
			"#name": "name",
		},
	}

	result, err := p.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	if len(result.Items) == 0 {
		return models.Products{}, nil
	}

	var product ProductTable
	products := make(models.Products, len(result.Items))
	for i, item := range result.Items {
		if err = attributevalue.UnmarshalMap(item, &product); err != nil {
			return nil, fmt.Errorf("error unmarshaling item: %w", err)
		}
		products[i] = assembleProduct(product)
	}

	return products, nil
}

func (p *Product) GetTopProducts(ctx context.Context) (models.Products, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String(p.tableName),
		FilterExpression: aws.String("is_top = :is_top"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":is_top": &types.AttributeValueMemberBOOL{Value: true},
		},
	}

	result, err := p.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	if len(result.Items) == 0 {
		return models.Products{}, nil
	}

	var product ProductTable
	products := make(models.Products, len(result.Items))
	for i, item := range result.Items {
		if err = attributevalue.UnmarshalMap(item, &product); err != nil {
			return nil, fmt.Errorf("error unmarshaling item: %w", err)
		}
		products[i] = assembleProduct(product)
	}

	return products, nil
}

func (p *Product) AddProduct(ctx context.Context, params mtypes.ProductParams) (*models.Product, error) {
	product := ProductTable{
		Uuid:         params.Uuid,
		Name:         params.Name,
		Price:        params.Price,
		Image:        params.Location,
		QRCode:       params.QRCode,
		IsTop:        params.IsTop,
		CategoryUuid: params.Category.Uuid,
		CategoryName: params.Category.Name,
		CreatedAt:    params.CreatedAt.Format(time.DateTime),
		UpdatedAt:    params.UpdatedAt.Format(time.DateTime),
	}

	item, err := attributevalue.MarshalMap(product)
	if err != nil {
		return nil, fmt.Errorf("error marshaling item: %w", err)
	}

	_, err = p.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(p.tableName),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("error adding item: %w", err)
	}

	return assembleProduct(product), nil
}

func (p *Product) PutProduct(ctx context.Context, params mtypes.ProductParams) (*models.Product, error) {
	expression := "SET #name = :name, price = :price, qrcode = :qrcode, is_top = :is_top, category_uuid = :category_uuid, category_name = :category_name, updated_at = :updated_at"
	expressionAttributeValues := map[string]types.AttributeValue{
		":name":          &types.AttributeValueMemberS{Value: params.Name},
		":price":         &types.AttributeValueMemberN{Value: strconv.FormatFloat(params.Price, 'f', -1, 32)},
		":qrcode":        &types.AttributeValueMemberS{Value: params.QRCode},
		":is_top":        &types.AttributeValueMemberBOOL{Value: params.IsTop},
		":category_uuid": &types.AttributeValueMemberS{Value: params.Category.Uuid},
		":category_name": &types.AttributeValueMemberS{Value: params.Category.Name},
		":updated_at":    &types.AttributeValueMemberS{Value: params.UpdatedAt.Format(time.DateTime)},
	}

	if params.Location != "" {
		expression += ", image = :image"
		expressionAttributeValues[":image"] = &types.AttributeValueMemberS{Value: params.Location}
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(p.tableName),
		Key: map[string]types.AttributeValue{
			"uuid": &types.AttributeValueMemberS{Value: params.Uuid},
		},
		UpdateExpression:          aws.String(expression),
		ConditionExpression:       aws.String("attribute_exists(#uuid)"),
		ExpressionAttributeValues: expressionAttributeValues,
		ExpressionAttributeNames: map[string]string{
			"#name": "name",
			"#uuid": "uuid",
		},
		ReturnValues: types.ReturnValueAllNew,
	}

	result, err := p.client.UpdateItem(ctx, input)
	if err != nil {
		var errf *types.ConditionalCheckFailedException
		if errors.As(err, &errf) {
			return nil, mtypes.ErrNotFound
		}

		return nil, fmt.Errorf("error updating item: %w", err)
	}

	var product ProductTable
	if err = attributevalue.UnmarshalMap(result.Attributes, &product); err != nil {
		return nil, fmt.Errorf("error marshaling item: %w", err)
	}

	return assembleProduct(product), nil
}

func (p *Product) DelProduct(ctx context.Context, uuid string) (*models.Product, error) {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(p.tableName),
		Key: map[string]types.AttributeValue{
			"uuid": &types.AttributeValueMemberS{Value: uuid},
		},
		ConditionExpression: aws.String("attribute_exists(#uuid)"),
		ExpressionAttributeNames: map[string]string{
			"#uuid": "uuid",
		},
		ReturnValues: types.ReturnValueAllOld,
	}

	result, err := p.client.DeleteItem(ctx, input)
	if err != nil {
		var errf *types.ConditionalCheckFailedException
		if errors.As(err, &errf) {
			return nil, mtypes.ErrNotFound
		}

		return nil, fmt.Errorf("error deleting item: %w", err)
	}

	var product ProductTable
	if err = attributevalue.UnmarshalMap(result.Attributes, &product); err != nil {
		return nil, fmt.Errorf("error marshaling item: %w", err)
	}

	return assembleProduct(product), nil
}

func assembleProduct(product ProductTable) *models.Product {
	return &models.Product{
		Uuid:   product.Uuid,
		Name:   product.Name,
		Price:  product.Price,
		Image:  product.Image,
		QRCode: product.QRCode,
		IsTop:  product.IsTop,
		Category: models.Category{
			Uuid: product.CategoryUuid,
			Name: product.CategoryName,
		},
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

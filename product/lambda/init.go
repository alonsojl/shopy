package main

import (
	"log"
	"log/slog"
	"os"
	"product/internal/apigateway"
	"product/internal/dynamodb"
	"product/internal/s3"
	"product/internal/service"
)

var handler *apigateway.Product

func init() {
	dynamoClient, err := dynamodb.Connection()
	if err != nil {
		log.Fatalf("error connecting to dynamodb: %v", err)
	}

	s3Client, err := s3.Connection()
	if err != nil {
		log.Fatalf("error connecting to amazon S3: %v", err)
	}

	var (
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		}))
		repository = dynamodb.NewProduct(logger, dynamoClient)
		storage    = s3.NewProduct(logger, s3Client)
		service    = service.NewProduct(logger, repository, storage)
	)

	handler = apigateway.NewProduct(logger, service)
}

package main

import (
	"category/internal/apigateway"
	"category/internal/dynamodb"
	"category/internal/s3"
	"category/internal/service"
	"log"
	"log/slog"
	"os"
)

var handler *apigateway.Category

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
		repository = dynamodb.NewCategory(logger, dynamoClient)
		storage    = s3.NewCategory(logger, s3Client)
		service    = service.NewCategory(logger, repository, storage)
	)

	handler = apigateway.NewCategory(logger, service)
}

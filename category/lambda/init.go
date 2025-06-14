package main

import (
	"log"
	"log/slog"
	"os"
	"shopy/internal/apigateway"
	"shopy/internal/dynamodb"
	"shopy/internal/s3"
	"shopy/internal/service"
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

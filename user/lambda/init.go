package main

import (
	"log"
	"log/slog"
	"os"
	"shopy/internal/apigateway"
	"shopy/internal/dynamodb"
	"shopy/internal/service"
)

var handler *apigateway.User

func init() {
	dynamoClient, err := dynamodb.Connection()
	if err != nil {
		log.Fatalf("error connecting to dynamodb: %v", err)
	}

	var (
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
		}))
		repository = dynamodb.NewUser(logger, dynamoClient)
		service    = service.NewUser(logger, repository)
	)

	handler = apigateway.NewUser(logger, service)
}

package apigateway

import (
	"category/internal/types"
	"category/pkg/errorx"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type APIGatewayFunc func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type ErrorResponse struct {
	BaseResponse
	Message string            `json:"message"`
	Errors  validation.Errors `json:"errors,omitempty" swaggertype:"object"`
}

func Error(err error) (events.APIGatewayProxyResponse, error) {
	var (
		errx     errorx.Error
		errs     validation.Errors
		response ErrorResponse
	)

	response.Datetime = time.Now().Format(time.DateTime)
	if !errors.As(err, &errx) {
		response.Status = "error"
		response.Message = "internal server error"
		response.Code = http.StatusInternalServerError
	} else {
		response.Status = "fail"
		response.Message = errx.Error()
		switch errx.Code() {
		case types.CodeInvalidArgument:
			response.Code = http.StatusBadRequest
			if errors.As(errx, &errs) {
				response.Errors = errs
			}
		case types.CodePrecondition:
			response.Code = http.StatusPreconditionFailed
		case types.CodeNotFound:
			response.Code = http.StatusNotFound
		}
	}
	return JSON(response, response.Code)
}

func JSON(response interface{}, statusCode int) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

package apigateway

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"shopy/internal/models"
	"shopy/internal/types"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type Service interface {
	GetCategories(ctx context.Context) (models.Categories, error)
	AddCategory(ctx context.Context, params types.CategoryParams) (*models.Category, error)
	DelCategory(ctx context.Context, uuid string) error
}

type Category struct {
	logger  *slog.Logger
	service Service
}

func NewCategory(logger *slog.Logger, service Service) *Category {
	return &Category{
		logger:  logger,
		service: service,
	}
}

func (c *Category) Router() APIGatewayFunc {
	return func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch event.HTTPMethod {
		case http.MethodGet:
			return c.HandleGetCategories(ctx, event)
		case http.MethodPost:
			return c.HandleAddCategory(ctx, event)
		case http.MethodDelete:
			return c.HandleDelCategory(ctx, event)
		}
		return events.APIGatewayProxyResponse{
			Body:       "method is not valid",
			StatusCode: http.StatusMethodNotAllowed,
		}, nil
	}
}

// @Summary 	Get categories.
// @Description Get product categories by store.
// @Tags 		Categories
// @Router 		/categories [get]
// @Accept 		json
// @Produce 	json
// @Security    JWT
// @Success     200	{object} SelectedCategories "Success"
// @Failure     400	{object} ErrorResponse "Bad Request"
// @Failure     401	{object} ErrorResponse "Unauthorized"
// @Failure     500	{object} ErrorResponse "Error Internal Server"
func (c *Category) HandleGetCategories(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	categories, err := c.service.GetCategories(ctx)
	if err != nil {
		c.logger.Error("error getting categories", "error", err)
		return Error(err)
	}

	var response = SelectedCategories{
		BaseResponse: NewBaseResponse(http.StatusOK),
		Categories:   categories,
	}

	return JSON(response, http.StatusOK)
}

// @Summary 	Add category.
// @Description Add new product category.
// @Tags 		Categories
// @Router 		/categories [post]
// @Accept 		json
// @Produce 	json
// @Security    JWT
// @Param	    params body  CategoryAddRequest true "Category"
// @Success     201	{object} CategoryAdded "Success"
// @Failure     400	{object} ErrorResponse "Bad Request"
// @Failure     401	{object} ErrorResponse "Unauthorized"
// @Failure     500	{object} ErrorResponse "Error Internal Server"
func (c *Category) HandleAddCategory(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var request CategoryAddRequest

	if err := json.Unmarshal([]byte(event.Body), &request); err != nil {
		c.logger.Error("invalid category body", "error", err)
		return Error(types.ErrRequest)
	}

	if err := request.Validate(); err != nil {
		c.logger.Error("invalid category params", "error", err)
		return Error(types.ErrParams.Wrap(err))
	}

	image, err := base64.StdEncoding.DecodeString(request.Image)
	if err != nil {
		c.logger.Error("error decoding image", "error", err)
		return Error(types.ErrRequest)
	}

	now := time.Now().UTC()
	category, err := c.service.AddCategory(ctx, types.CategoryParams{
		Uuid:      uuid.New().String(),
		Name:      request.Name,
		Image:     image,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		c.logger.Error("error adding category", "error", err)
		return Error(err)
	}

	var response = CategoryAdded{
		BaseResponse: NewBaseResponse(http.StatusCreated),
		Category:     category,
	}

	return JSON(response, http.StatusCreated)
}

// @Summary 	Delete category.
// @Description Delete product category.
// @Tags 		Categories
// @Router 		/categories/{uuid} [delete]
// @Accept 		json
// @Produce 	json
// @Security    JWT
// @Param       uuid path string true "Category UUID"
// @Success     200	{object} CategoryDeleted "Success"
// @Failure     400	{object} ErrorResponse "Bad Request"
// @Failure     401	{object} ErrorResponse "Unauthorized"
// @Failure     500	{object} ErrorResponse "Error Internal Server"
func (c *Category) HandleDelCategory(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	uuid := event.PathParameters["uuid"]
	if err := c.service.DelCategory(ctx, uuid); err != nil {
		c.logger.Error("error deleting category", "error", err)
		return Error(err)
	}

	var response = CategoryDeleted{
		BaseResponse: NewBaseResponse(http.StatusOK),
		Category:     "deleted",
	}

	return JSON(response, http.StatusOK)
}

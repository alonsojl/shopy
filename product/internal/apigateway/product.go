package apigateway

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"shopy/internal/domain"
	"shopy/internal/models"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type Service interface {
	SearchProducts(ctx context.Context, params domain.ProductParams) (models.Products, error)
	AddProduct(ctx context.Context, params domain.ProductParams) (*models.Product, error)
	PutProduct(ctx context.Context, params domain.ProductParams) (*models.Product, error)
	DelProduct(ctx context.Context, uuid string) error
}

type Product struct {
	logger  *slog.Logger
	service Service
}

func NewProduct(logger *slog.Logger, service Service) *Product {
	return &Product{
		logger:  logger,
		service: service,
	}
}

func (p *Product) Router() APIGatewayFunc {
	return func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch event.HTTPMethod {
		case http.MethodGet:
			return p.HandleSearchProducts(ctx, event)
		case http.MethodPost:
			return p.HandleAddProduct(ctx, event)
		case http.MethodPut:
			return p.HandlePutProduct(ctx, event)
		case http.MethodDelete:
			return p.HandleDelProduct(ctx, event)
		}
		return events.APIGatewayProxyResponse{
			Body:       "method is not valid",
			StatusCode: http.StatusMethodNotAllowed,
		}, nil
	}
}

// @Summary 	Get products.
// @Description Retrieves products using query parameters, only one parameter can be used at a time. By default, the top products are returned.
// @Tags 		Products
// @Router 		/products [get]
// @Accept 		json
// @Produce 	json
// @Security    JWT
// @Param       name query string true "Product name"
// @Param       qrcode query string true "Product QR code"
// @Param       category_uuid query string true "Product category UUID"
// @Success     200	{object} SelectedProducts "Success"
// @Failure     400	{object} ErrorResponse "Bad Request"
// @Failure     401	{object} ErrorResponse "Unauthorized"
// @Failure     500	{object} ErrorResponse "Error Internal Server"
func (p *Product) HandleSearchProducts(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	products, err := p.service.SearchProducts(ctx, domain.ProductParams{
		Name:   event.QueryStringParameters["name"],
		QRCode: event.QueryStringParameters["qrcode"],
		Category: domain.Category{
			Uuid: event.QueryStringParameters["category_uuid"],
		},
	})
	if err != nil {
		p.logger.Error("error getting products", "error", err)
		return Error(err)
	}

	var response = SelectedProducts{
		BaseResponse: NewBaseResponse(http.StatusOK),
		Products:     products,
	}

	return JSON(response, http.StatusOK)
}

// @Summary 	Add product.
// @Description Add new product.
// @Tags 		Products
// @Router 		/products [post]
// @Accept 		json
// @Produce 	json
// @Security    JWT
// @Param	    params body  ProductAddRequest true "Product"
// @Success     201	{object} ProductAdded "Success"
// @Failure     400	{object} ErrorResponse "Bad Request"
// @Failure     401	{object} ErrorResponse "Unauthorized"
// @Failure     500	{object} ErrorResponse "Error Internal Server"
func (p *Product) HandleAddProduct(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var request ProductAddRequest

	if err := json.Unmarshal([]byte(event.Body), &request); err != nil {
		p.logger.Error("invalid product body", "error", err)
		return Error(domain.ErrRequest)
	}

	if err := request.Validate(); err != nil {
		p.logger.Error("invalid product params", "error", err)
		return Error(domain.ErrParams.Wrap(err))
	}

	image, err := base64.StdEncoding.DecodeString(request.Image)
	if err != nil {
		p.logger.Error("error decoding image", "error", err)
		return Error(domain.ErrRequest)
	}

	now := time.Now().UTC()
	product, err := p.service.AddProduct(ctx, domain.ProductParams{
		Uuid:   uuid.New().String(),
		Name:   request.Name,
		Price:  request.Price,
		QRCode: request.QRCode,
		IsTop:  request.IsTop,
		Category: domain.Category{
			Uuid: request.Category.Uuid,
			Name: request.Category.Name,
		},
		Image:     image,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		p.logger.Error("error adding product", "error", err)
		return Error(err)
	}

	var response = ProductAdded{
		BaseResponse: NewBaseResponse(http.StatusCreated),
		Product:      product,
	}

	return JSON(response, http.StatusCreated)
}

// @Summary 	Update product.
// @Description Update product and image fields.
// @Tags 		Products
// @Router 		/products/{uuid} [put]
// @Accept 		json
// @Produce 	json
// @Security    JWT
// @Param       uuid path string true "Product UUID"
// @Param	    params body  ProductPutRequest true "Product"
// @Success     201	{object} ProductAdded "Success"
// @Failure     400	{object} ErrorResponse "Bad Request"
// @Failure     401	{object} ErrorResponse "Unauthorized"
// @Failure     500	{object} ErrorResponse "Error Internal Server"
func (p *Product) HandlePutProduct(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		err     error
		request ProductPutRequest
	)

	if err = json.Unmarshal([]byte(event.Body), &request); err != nil {
		p.logger.Error("invalid product body", "error", err)
		return Error(domain.ErrRequest)
	}

	if err = request.Validate(); err != nil {
		p.logger.Error("invalid product params", "error", err)
		return Error(domain.ErrParams.Wrap(err))
	}

	var image []byte
	if request.Image != "" {
		image, err = base64.StdEncoding.DecodeString(request.Image)
		if err != nil {
			p.logger.Error("error decoding image", "error", err)
			return Error(domain.ErrRequest)
		}
	}

	uuid := event.PathParameters["uuid"]
	product, err := p.service.PutProduct(ctx, domain.ProductParams{
		Uuid:   uuid,
		Name:   request.Name,
		Price:  request.Price,
		QRCode: request.QRCode,
		IsTop:  request.IsTop,
		Category: domain.Category{
			Uuid: request.Category.Uuid,
			Name: request.Category.Name,
		},
		Image:     image,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		p.logger.Error("error updating product", "error", err)
		return Error(err)
	}

	var response = ProductAdded{
		BaseResponse: NewBaseResponse(http.StatusOK),
		Product:      product,
	}

	return JSON(response, http.StatusOK)
}

// @Summary 	Delete product.
// @Description Delete product and image.
// @Tags 		Products
// @Router 		/products/{uuid} [delete]
// @Accept 		json
// @Produce 	json
// @Security    JWT
// @Param       uuid path string true "Product UUID"
// @Success     200	{object} ProductDeleted "Success"
// @Failure     400	{object} ErrorResponse "Bad Request"
// @Failure     401	{object} ErrorResponse "Unauthorized"
// @Failure     500	{object} ErrorResponse "Error Internal Server"
func (p *Product) HandleDelProduct(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	uuid := event.PathParameters["uuid"]
	if err := p.service.DelProduct(ctx, uuid); err != nil {
		p.logger.Error("error deleting product", "error", err)
		return Error(err)
	}

	var response = ProductDeleted{
		BaseResponse: NewBaseResponse(http.StatusOK),
		Product:      "deleted",
	}

	return JSON(response, http.StatusOK)
}

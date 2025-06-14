package service

import (
	"context"
	"log/slog"
	"shopy/internal/models"
	"shopy/internal/types"
)

type Repository interface {
	GetProductsByCategory(ctx context.Context, uuid string) (models.Products, error)
	GetProductsByQRCode(ctx context.Context, qrcode string) (models.Products, error)
	GetProductsByName(ctx context.Context, name string) (models.Products, error)
	GetTopProducts(ctx context.Context) (models.Products, error)
	AddProduct(ctx context.Context, params types.ProductParams) (*models.Product, error)
	PutProduct(ctx context.Context, params types.ProductParams) (*models.Product, error)
	DelProduct(ctx context.Context, uuid string) (*models.Product, error)
}

type Storage interface {
	UploadImage(ctx context.Context, uuid string, image []byte) (string, error)
	DeleteImage(ctx context.Context, image string) error
}

type Product struct {
	logger     *slog.Logger
	repository Repository
	storage    Storage
}

func NewProduct(logger *slog.Logger, repository Repository, storage Storage) *Product {
	return &Product{
		logger:     logger,
		repository: repository,
		storage:    storage,
	}
}

func (p *Product) SearchProducts(ctx context.Context, params types.ProductParams) (models.Products, error) {
	switch {
	case params.Category.Uuid != "":
		return p.repository.GetProductsByCategory(ctx, params.Category.Uuid)
	case params.QRCode != "":
		return p.repository.GetProductsByQRCode(ctx, params.QRCode)
	case params.Name != "":
		return p.repository.GetProductsByName(ctx, params.Name)
	default:
		return p.repository.GetTopProducts(ctx)
	}
}

func (p *Product) AddProduct(ctx context.Context, params types.ProductParams) (*models.Product, error) {
	location, err := p.storage.UploadImage(ctx, params.Uuid, params.Image)
	if err != nil {
		return nil, err
	}

	params.Location = location
	return p.repository.AddProduct(ctx, params)
}

func (p *Product) PutProduct(ctx context.Context, params types.ProductParams) (*models.Product, error) {
	if params.Image != nil {
		location, err := p.storage.UploadImage(ctx, params.Uuid, params.Image)
		if err != nil {
			return nil, err
		}
		params.Location = location
	}

	return p.repository.PutProduct(ctx, params)
}

func (p *Product) DelProduct(ctx context.Context, uuid string) error {
	product, err := p.repository.DelProduct(ctx, uuid)
	if err != nil {
		return err
	}
	return p.storage.DeleteImage(ctx, product.Image)
}

package service

import (
	"context"
	"log/slog"
	"shopy/internal/domain"
	"shopy/internal/models"
)

type Repository interface {
	GetCategories(ctx context.Context) (models.Categories, error)
	AddCategory(ctx context.Context, params domain.CategoryParams) (*models.Category, error)
	DelCategory(ctx context.Context, uuid string) (*models.Category, error)
}

type Storage interface {
	UploadImage(ctx context.Context, uuid string, image []byte) (string, error)
	DeleteImage(ctx context.Context, image string) error
}

type Category struct {
	logger     *slog.Logger
	repository Repository
	storage    Storage
}

func NewCategory(logger *slog.Logger, repository Repository, storage Storage) *Category {
	return &Category{
		logger:     logger,
		repository: repository,
		storage:    storage,
	}
}

func (c *Category) GetCategories(ctx context.Context) (models.Categories, error) {
	return c.repository.GetCategories(ctx)
}

func (c *Category) AddCategory(ctx context.Context, params domain.CategoryParams) (*models.Category, error) {
	location, err := c.storage.UploadImage(ctx, params.Uuid, params.Image)
	if err != nil {
		return nil, err
	}

	params.Location = location
	category, err := c.repository.AddCategory(ctx, params)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (c *Category) DelCategory(ctx context.Context, uuid string) error {
	category, err := c.repository.DelCategory(ctx, uuid)
	if err != nil {
		return err
	}
	return c.storage.DeleteImage(ctx, category.Image)
}

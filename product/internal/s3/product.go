package s3

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const folderName = "product"

type Product struct {
	logger *slog.Logger
	client *s3.Client
	bucket string
}

func NewProduct(logger *slog.Logger, client *s3.Client) *Product {
	bucket := os.Getenv("BUCKET_NAME")
	return &Product{
		logger: logger,
		client: client,
		bucket: bucket,
	}
}

func (p *Product) UploadImage(ctx context.Context, uuid string, image []byte) (string, error) {
	var (
		contentType = http.DetectContentType(image)
		filename    = fmt.Sprintf("%s.%s", uuid, "png")
		filepath    = fmt.Sprintf("%s/%s", folderName, filename)
		uploader    = manager.NewUploader(p.client)
	)

	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(filepath),
		Body:        bytes.NewReader(image),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("error uploading image: %w", err)
	}

	return result.Location, nil
}

func (p *Product) DeleteImage(ctx context.Context, location string) error {
	var (
		filename = path.Base(location)
		filepath = fmt.Sprintf("%s/%s", folderName, filename)
	)
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(filepath),
	})
	if err != nil {
		return fmt.Errorf("error deleting image: %w", err)
	}

	return nil
}

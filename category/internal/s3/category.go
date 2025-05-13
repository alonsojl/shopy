package s3

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const folderName = "category"

type Category struct {
	logger *slog.Logger
	client *s3.Client
	bucket string
}

func NewCategory(logger *slog.Logger, client *s3.Client) *Category {
	bucket := os.Getenv("BUCKET_NAME")
	return &Category{
		logger: logger,
		client: client,
		bucket: bucket,
	}
}

func (c *Category) UploadImage(ctx context.Context, uuid string, image []byte) (string, error) {
	var (
		contentType = http.DetectContentType(image)
		extension   = strings.Split(contentType, "/")[1]
		filename    = fmt.Sprintf("%s.%s", uuid, extension)
		filepath    = fmt.Sprintf("%s/%s", folderName, filename)
		uploader    = manager.NewUploader(c.client)
	)

	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(filepath),
		Body:        bytes.NewReader(image),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("error uploading image: %w", err)
	}

	return result.Location, nil
}

func (c *Category) DeleteImage(ctx context.Context, location string) error {
	var (
		filename = path.Base(location)
		filepath = fmt.Sprintf("%s/%s", folderName, filename)
	)
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(filepath),
	})
	if err != nil {
		return fmt.Errorf("error deleting image: %w", err)
	}

	return nil
}

package bucket

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

type Bucket struct {
	*s3.Client
	BucketName string
}

func NewBucket(cfg aws.Config, bucketName string) *Bucket {
	return &Bucket{s3.NewFromConfig(cfg), bucketName}
}

func (b *Bucket) SendFileToS3(ctx context.Context, file models.FileData) (string, error) {
	inputPut := &s3.PutObjectInput{
		Bucket:      aws.String(b.BucketName),
		Key:         aws.String(file.Filename),
		Body:        strings.NewReader(string(file.Content)),
		ContentType: aws.String(file.ContentType),
	}

	_, err := b.PutObject(ctx, inputPut)
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", b.BucketName, file.Filename), nil
}

func (b *Bucket) FileExistsInS3(ctx context.Context, fileName string) (bool, error) {
	_, err := b.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(b.BucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		// Check if it's a "not found" error
		var nf *types.NoSuchKey
		if errors.As(err, &nf) {
			return false, nil // File doesn't exist, but no error
		}
		return false, err // Some other error occurred
	}
	return true, nil
}

func (b *Bucket) DeleteFileFromS3(ctx context.Context, fileName string) error {
	inputDelete := &s3.DeleteObjectInput{
		Bucket: aws.String(b.BucketName),
		Key:    aws.String(fileName),
	}

	_, err := b.DeleteObject(ctx, inputDelete)
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

func (b *Bucket) GeneratePresignedURL(ctx context.Context, fileName string) (*v4.PresignedHTTPRequest, error) {
	inputGet := &s3.GetObjectInput{
		Bucket: aws.String(b.BucketName),
		Key:    aws.String(fileName),
	}

	presignClient := s3.NewPresignClient(b.Client)
	presignedReq, err := presignClient.PresignGetObject(ctx, inputGet, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(60 * time.Minute) // URL expires in 1 hour
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedReq, nil
}

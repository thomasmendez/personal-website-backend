package bucket

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

type Bucket struct {
	*s3.Client
}

func NewBucket(cfg aws.Config) *Bucket {
	return &Bucket{s3.NewFromConfig(cfg)}
}

func SendFileToS3(ctx context.Context, svc *s3.Client, bucketName string, file models.FileData) (string, error) {
	inputPut := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(file.Filename),
		Body:        strings.NewReader(string(file.Content)),
		ContentType: aws.String(file.ContentType),
	}

	_, err := svc.PutObject(ctx, inputPut)
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, file.Filename), nil
}

func DeleteFileFromS3(ctx context.Context, svc *s3.Client, bucketName string, mediaLink string) error {
	fileName, err := GetFileNameFromMediaLink(mediaLink)
	if err != nil {
		return fmt.Errorf("failed to get filename from mediaLink: %w", err)
	}

	inputDelete := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	}

	_, err = svc.DeleteObject(ctx, inputDelete)
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

func GeneratePresignedURL(ctx context.Context, svc *s3.Client, bucketName string, mediaLink string) (*v4.PresignedHTTPRequest, error) {
	fileName, err := GetFileNameFromMediaLink(mediaLink)
	if err != nil {
		return nil, fmt.Errorf("failed to get filename from mediaLink: %w", err)
	}

	inputGet := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	}

	presignClient := s3.NewPresignClient(svc)
	presignedReq, err := presignClient.PresignGetObject(ctx, inputGet, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(60 * time.Minute) // URL expires in 1 hour
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedReq, nil
}

func GetFileNameFromMediaLink(mediaLink string) (string, error) {
	parsedURL, err := url.Parse(mediaLink)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}

	filename := path.Base(parsedURL.Path)

	// Check if we actually got a filename (not empty or just "/")
	if filename == "." || filename == "/" || filename == "" {
		return "", fmt.Errorf("no filename found in MediaLink")
	}

	return filename, nil
}

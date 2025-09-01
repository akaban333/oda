package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/studyplatform/backend/pkg/logger"
)

// MinioClient represents a MinIO client for file storage
type MinioClient struct {
	client      *minio.Client
	buckets     map[string]string
	useSSL      bool
	expiry      time.Duration
	initialized bool
}

// NewMinioClient creates a new MinIO client for file storage
func NewMinioClient() (*MinioClient, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSLStr := os.Getenv("MINIO_USE_SSL")

	// Set default values if not provided
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	if accessKey == "" {
		accessKey = "minioadmin"
	}

	if secretKey == "" {
		secretKey = "minioadmin"
	}

	useSSL := false
	if useSSLStr != "" {
		var err error
		useSSL, err = strconv.ParseBool(useSSLStr)
		if err != nil {
			useSSL = false
		}
	}

	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating MinIO client: %w", err)
	}

	// Create bucket map
	buckets := map[string]string{
		"avatars":   os.Getenv("MINIO_AVATAR_BUCKET"),
		"materials": os.Getenv("MINIO_MATERIALS_BUCKET"),
		"temp":      os.Getenv("MINIO_TEMP_BUCKET"),
	}

	// Set default values for buckets if not provided
	if buckets["avatars"] == "" {
		buckets["avatars"] = "avatars"
	}

	if buckets["materials"] == "" {
		buckets["materials"] = "materials"
	}

	if buckets["temp"] == "" {
		buckets["temp"] = "temp"
	}

	return &MinioClient{
		client:      minioClient,
		buckets:     buckets,
		useSSL:      useSSL,
		expiry:      time.Hour * 24, // Default expiry for pre-signed URLs
		initialized: false,
	}, nil
}

// Initialize creates the required buckets if they don't exist
func (m *MinioClient) Initialize(ctx context.Context) error {
	if m.initialized {
		return nil
	}

	// Ensure all buckets exist
	for _, bucketName := range m.buckets {
		exists, err := m.client.BucketExists(ctx, bucketName)
		if err != nil {
			return fmt.Errorf("error checking if bucket %s exists: %w", bucketName, err)
		}

		if !exists {
			err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
			if err != nil {
				return fmt.Errorf("error creating bucket %s: %w", bucketName, err)
			}
			logger.Info(fmt.Sprintf("Created bucket %s", bucketName))
		}
	}

	m.initialized = true
	return nil
}

// UploadFile uploads a file to MinIO
func (m *MinioClient) UploadFile(ctx context.Context, bucketType string, objectName string, reader io.Reader, size int64, contentType string) error {
	if !m.initialized {
		if err := m.Initialize(ctx); err != nil {
			return err
		}
	}

	bucketName, ok := m.buckets[bucketType]
	if !ok {
		return fmt.Errorf("invalid bucket type: %s", bucketType)
	}

	_, err := m.client.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("error uploading file to MinIO: %w", err)
	}

	return nil
}

// DownloadFile downloads a file from MinIO
func (m *MinioClient) DownloadFile(ctx context.Context, bucketType string, objectName string) (io.ReadCloser, error) {
	if !m.initialized {
		if err := m.Initialize(ctx); err != nil {
			return nil, err
		}
	}

	bucketName, ok := m.buckets[bucketType]
	if !ok {
		return nil, fmt.Errorf("invalid bucket type: %s", bucketType)
	}

	obj, err := m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("error downloading file from MinIO: %w", err)
	}

	return obj, nil
}

// GetPresignedURL generates a pre-signed URL for object
func (m *MinioClient) GetPresignedURL(ctx context.Context, bucketType string, objectName string, method string) (string, error) {
	if !m.initialized {
		if err := m.Initialize(ctx); err != nil {
			return "", err
		}
	}

	bucketName, ok := m.buckets[bucketType]
	if !ok {
		return "", fmt.Errorf("invalid bucket type: %s", bucketType)
	}

	// Get presigned URL based on method
	if method == "GET" {
		presignedURL, err := m.client.PresignedGetObject(ctx, bucketName, objectName, m.expiry, url.Values{})
		if err != nil {
			return "", fmt.Errorf("error generating presigned GET URL: %w", err)
		}
		return presignedURL.String(), nil
	} else if method == "PUT" {
		presignedURL, err := m.client.PresignedPutObject(ctx, bucketName, objectName, m.expiry)
		if err != nil {
			return "", fmt.Errorf("error generating presigned PUT URL: %w", err)
		}
		return presignedURL.String(), nil
	}

	return "", fmt.Errorf("unsupported method: %s", method)
}

// DeleteFile deletes a file from MinIO
func (m *MinioClient) DeleteFile(ctx context.Context, bucketType string, objectName string) error {
	if !m.initialized {
		if err := m.Initialize(ctx); err != nil {
			return err
		}
	}

	bucketName, ok := m.buckets[bucketType]
	if !ok {
		return fmt.Errorf("invalid bucket type: %s", bucketType)
	}

	err := m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("error deleting file from MinIO: %w", err)
	}

	return nil
}

// ListFiles lists all files in a bucket
func (m *MinioClient) ListFiles(ctx context.Context, bucketType string, prefix string) ([]string, error) {
	if !m.initialized {
		if err := m.Initialize(ctx); err != nil {
			return nil, err
		}
	}

	bucketName, ok := m.buckets[bucketType]
	if !ok {
		return nil, fmt.Errorf("invalid bucket type: %s", bucketType)
	}

	objectCh := m.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var objectNames []string
	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("error listing objects: %w", object.Err)
		}
		objectNames = append(objectNames, object.Key)
	}

	return objectNames, nil
}

// GetPresignedUploadURL generates a presigned URL for file upload
func (m *MinioClient) GetPresignedUploadURL(bucketType string, objectName string, contentType string, expiry time.Duration) (string, error) {
	ctx := context.Background()
	if !m.initialized {
		if err := m.Initialize(ctx); err != nil {
			return "", err
		}
	}

	bucketName, ok := m.buckets[bucketType]
	if !ok {
		return "", fmt.Errorf("invalid bucket type: %s", bucketType)
	}

	// Set content type policy
	policy := url.Values{}
	if contentType != "" {
		policy.Set("Content-Type", contentType)
	}

	presignedURL, err := m.client.PresignedPutObject(ctx, bucketName, objectName, expiry)
	if err != nil {
		return "", fmt.Errorf("error generating presigned upload URL: %w", err)
	}

	return presignedURL.String(), nil
}

// FileExists checks if a file exists in the specified bucket
func (m *MinioClient) FileExists(bucketType string, objectName string) (bool, error) {
	ctx := context.Background()
	if !m.initialized {
		if err := m.Initialize(ctx); err != nil {
			return false, err
		}
	}

	bucketName, ok := m.buckets[bucketType]
	if !ok {
		return false, fmt.Errorf("invalid bucket type: %s", bucketType)
	}

	_, err := m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		// Check if error is because object doesn't exist
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("error checking file existence: %w", err)
	}

	return true, nil
}

// GetFileURL generates a public URL for accessing a file
func (m *MinioClient) GetFileURL(bucketType string, objectName string) (string, error) {
	ctx := context.Background()
	if !m.initialized {
		if err := m.Initialize(ctx); err != nil {
			return "", err
		}
	}

	bucketName, ok := m.buckets[bucketType]
	if !ok {
		return "", fmt.Errorf("invalid bucket type: %s", bucketType)
	}

	// Generate a presigned URL for GET with long expiry (7 days)
	presignedURL, err := m.client.PresignedGetObject(ctx, bucketName, objectName, 7*24*time.Hour, url.Values{})
	if err != nil {
		return "", fmt.Errorf("error generating file URL: %w", err)
	}

	return presignedURL.String(), nil
}

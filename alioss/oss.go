package alioss

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"path"
	"time"
)

type Client struct {
	cfg    Config
	client *oss.Client
	bucket *oss.Bucket
}

// New 创建 OSS 客户端
func New(cfg Config) (*Client, error) {
	if cfg.MaxRetry <= 0 {
		cfg.MaxRetry = 3
	}

	cli, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OSS client: %w", err)
	}

	bucket, err := cli.Bucket(cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}

	return &Client{
		cfg:    cfg,
		client: cli,
		bucket: bucket,
	}, nil
}

func (c *Client) tryUpload(key string, reader io.Reader) error {
	var err error
	for i := 0; i < c.cfg.MaxRetry; i++ {
		err = c.bucket.PutObject(key, reader)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return fmt.Errorf("upload failed after %d retries: %w", c.cfg.MaxRetry, err)
}

// UploadBytes uploads a file from byte slice and returns its URL.
func (c *Client) UploadBytes(fileBytes []byte, ext string) (string, error) {
	var key string
	if c.cfg.EnableHashCheck {
		hash := md5.Sum(fileBytes)
		hashStr := hex.EncodeToString(hash[:])
		key = path.Join(c.cfg.BasePath, hashStr+ext)
		exist, err := c.bucket.IsObjectExist(key)
		if err == nil && exist {
			return c.generateURL(key)
		}
	} else {
		fileName := uuid.New().String() + ext
		key = path.Join(c.cfg.BasePath, fileName)
	}

	err := c.tryUpload(key, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}

	return c.generateURL(key)
}

// UploadFromMultipart uploads from multipart.FileHeader (e.g., from HTTP upload)
func (c *Client) UploadFromMultipart(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("open multipart file error: %w", err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("read multipart file error: %w", err)
	}

	ext := path.Ext(fileHeader.Filename)
	return c.UploadBytes(fileBytes, ext)
}

// GetSignedURL returns a signed URL with expiration (only for private buckets)
func (c *Client) GetSignedURL(objectKey string, expireSeconds int64) (string, error) {
	key := path.Join(c.cfg.BasePath, objectKey)
	url, err := c.bucket.SignURL(key, oss.HTTPGet, expireSeconds)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	return url, nil
}

// Delete deletes an object from OSS
func (c *Client) Delete(objectKey string) error {
	key := path.Join(c.cfg.BasePath, objectKey)
	return c.bucket.DeleteObject(key)
}

// generateURL generates the public or signed URL based on config
func (c *Client) generateURL(objectKey string) (string, error) {
	if c.cfg.IsPrivate {
		return c.GetSignedURL(objectKey, 3600)
	}
	return fmt.Sprintf("https://%s.%s/%s", c.cfg.BucketName, c.cfg.Endpoint, objectKey), nil
}

// UploadFile uploads a local file to OSS and returns its URL
func (c *Client) UploadFile(localPath string) (string, error) {
	fileBytes, err := os.ReadFile(localPath)
	if err != nil {
		return "", fmt.Errorf("read local file error: %w", err)
	}
	ext := path.Ext(localPath)
	return c.UploadBytes(fileBytes, ext)
}

// UploadLargeFile uses multipart upload for large files
func (c *Client) UploadLargeFile(localPath string, partSize int64) (string, error) {
	fileName := uuid.New().String() + path.Ext(localPath)
	key := path.Join(c.cfg.BasePath, fileName)

	var err error
	for i := 0; i < c.cfg.MaxRetry; i++ {
		err = c.bucket.UploadFile(key, localPath, partSize, oss.Routines(3))
		if err == nil {
			return c.generateURL(key)
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return "", fmt.Errorf("multipart upload failed after %d retries: %w", c.cfg.MaxRetry, err)
}

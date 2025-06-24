package alioss

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"path"
)

type Client struct {
	cfg    Config
	client *oss.Client
	bucket *oss.Bucket
}

// New 创建 OSS 客户端
func New(cfg Config) (*Client, error) {
	c, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("init OSS client error: %w", err)
	}

	bucket, err := c.Bucket(cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("get OSS bucket error: %w", err)
	}

	return &Client{
		cfg:    cfg,
		client: c,
		bucket: bucket,
	}, nil
}

// Upload 上传文件（通过 []byte）
func (c *Client) Upload(objectName string, content []byte) error {
	key := path.Join(c.cfg.BasePath, objectName)
	return c.bucket.PutObject(key, bytes.NewReader(content))
}

// UploadFile 上传本地文件
func (c *Client) UploadFile(objectName, localFilePath string) error {
	key := path.Join(c.cfg.BasePath, objectName)
	return c.bucket.PutObjectFromFile(key, localFilePath)
}

// GetURL 获取文件完整访问链接（公开读）
func (c *Client) GetURL(objectName string) string {
	key := path.Join(c.cfg.BasePath, objectName)
	return fmt.Sprintf("https://%s.%s/%s", c.cfg.BucketName, c.cfg.Endpoint, key)
}

// GetSignedURL 获取带签名的私有访问链接（有效期：秒）
func (c *Client) GetSignedURL(objectName string, expireSeconds int64) (string, error) {
	key := path.Join(c.cfg.BasePath, objectName)
	signedURL, err := c.bucket.SignURL(key, oss.HTTPGet, expireSeconds)
	if err != nil {
		return "", fmt.Errorf("sign URL error: %w", err)
	}
	return signedURL, nil
}

// Delete 删除文件
func (c *Client) Delete(objectName string) error {
	key := path.Join(c.cfg.BasePath, objectName)
	return c.bucket.DeleteObject(key)
}

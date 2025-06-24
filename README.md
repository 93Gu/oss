README - Aliyun OSS Go SDK
==========================

一个简洁易用的阿里云 OSS 上传封装，支持公开/私有存储桶，支持多种上传方式。

依赖：
- github.com/aliyun/aliyun-oss-go-sdk/oss
- github.com/google/uuid

安装：
```
go get github.com/aliyun/aliyun-oss-go-sdk/oss
```

使用示例：
```go
import "path/to/internal/sdk/oss"

cfg := osssdk.Config{
	Endpoint:        "oss-cn-hangzhou.aliyuncs.com",
	AccessKeyID:     "your-access-key-id",
	AccessKeySecret: "your-access-key-secret",
	BucketName:      "your-bucket-name",
	BasePath:        "uploads",
	IsPrivate:       true,
}

client, err := osssdk.New(cfg)
if err != nil {
	panic(err)
}

// 上传本地文件
url, err := client.UploadFile("/tmp/demo.jpg")
fmt.Println("File uploaded URL:", url)

// 删除文件（传入上传时返回的文件名）
err = client.Delete("uploads/xxxx.jpg")

// 获取签名链接（私有桶）
signedURL, _ := client.GetSignedURL("uploads/xxxx.jpg", 3600)
```

支持功能：
- ✅ 上传本地文件
- ✅ 上传 HTTP 文件（`multipart.FileHeader`）
- ✅ 获取访问链接（公开 or 私有签名）
- ✅ 删除文件
- ✅ 自动生成唯一文件名（UUID）

未来计划：
- 分片上传
- 下载接口封装
- 上传失败重试
- 上传进度显示

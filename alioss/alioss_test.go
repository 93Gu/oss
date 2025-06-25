package alioss

import (
	"fmt"
	"testing"
)

func TestOSS(t *testing.T) {
	client, err := New(Config{
		Endpoint:        "oss-cn-hangzhou.aliyuncs.com",
		AccessKeyID:     "your-access-key-id",
		AccessKeySecret: "your-access-key-secret",
		BucketName:      "your-bucket-name",
		BasePath:        "prefix", // 存储路径前缀，可为空
	})
	if err != nil {
		t.Fatal(err)
	}

	// 上传文件
	url, err := client.UploadFile("./hello.txt")
	if err != nil {
		t.Fatal("upload error:", err)
	}

	fmt.Println("Public URL:", url)

	// 获取签名链接
	signedURL, _ := client.GetSignedURL("hello.txt", 3600)
	fmt.Println("Signed URL:", signedURL)

	// 删除文件
	// _ = client.Delete("hello.txt")
}

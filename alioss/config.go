package alioss

type Config struct {
	Endpoint        string // 如：https://oss-cn-hangzhou.aliyuncs.com
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	BasePath        string // 可选：指定文件统一前缀路径
	IsPrivate       bool   // true = use signed URL, false = public
	EnableHashCheck bool   // true = 启用秒传（基于 md5）
	MaxRetry        int    // 上传失败最大重试次数，默认可设为 3
}

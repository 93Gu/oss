package alioss

type Config struct {
	Endpoint        string // 如：https://oss-cn-hangzhou.aliyuncs.com
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	BasePath        string // 可选：指定文件统一前缀路径
}

package alioss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// Download 下载oss文件
// endpoint: oss的接入地址，例如oss-ap-southeast-1.aliyuncs.com。 通过oss-进入具体的bucket-概览页面查看。
// bucketName: 桶名称
// ossFilePath: 您要下载的文件位于桶中的路径
// localPathPrefix: 希望保存的文件名(不带扩展名)，扩展名在返回值localPath中拼接。
// 返回: 如果失败，返回err结果。 如果成功，返回文件下载后的路径地址
func Download(endpoint string, ak string, sk string, bucketName string, ossFilePath string, localPathPrefix string) (localPath string, err error) {
	client, err := oss.New(endpoint, ak, sk) // 创建OSSClient实例
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", err
	}
	props, err := bucket.GetObjectDetailedMeta(ossFilePath)
	if err != nil {
		return "", err
	}

	// 文件后缀. 文件类型在HTTP头"Content-Type"中
	var extension string
	switch props.Get("Content-Type") {
	case "application/x-zip-compressed":
		extension = ".zip"
	case "application/csv":
		extension = ".csv"
	default:
		extension = ""
	}
	localPath = fmt.Sprintf("%s%s", localPathPrefix, extension)
	err = bucket.GetObjectToFile(ossFilePath, localPath)
	if err != nil {
		return "", err
	}
	return localPath, nil
}

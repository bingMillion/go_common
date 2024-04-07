package alioss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// Download 下载oss文件
// endpoint: oss的接入地址，例如oss-ap-southeast-1.aliyuncs.com。 通过oss-进入具体的bucket-概览页面查看。
// bucketName: 桶名称
// ossFilePath: 您要下载的文件位于桶中的路径
// localPath: 希望保存的文件名(不带扩展名) .如2023-01-ali
// 返回: 如果失败，返回err结果。 如果成功，isZip表示文件是否为zip格式
func Download(endpoint string, ak string, sk string, bucketName string, ossFilePath string, localPath string) (isZip bool, err error) {
	client, err := oss.New(endpoint, ak, sk) // 创建OSSClient实例
	if err != nil {
		return false, err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return false, err
	}
	props, err := bucket.GetObjectDetailedMeta(ossFilePath)
	if err != nil {
		return false, err
	}

	// 文件后缀. 文件类型在HTTP头"Content-Type"中
	var extension string
	switch props.Get("Content-Type") {
	case "application/x-zip-compressed":
		extension = ".zip"
		isZip = true
	case "application/csv":
		extension = ".csv"
	default:
		extension = ""
	}
	err = bucket.GetObjectToFile(ossFilePath, fmt.Sprintf("%s%s", localPath, extension))
	if err != nil {
		return false, err
	}
	return isZip, nil
}

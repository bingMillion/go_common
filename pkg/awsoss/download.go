package awsoss

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

// Download
// bucketName: 桶名称
// region 区域，例如 us-west-2
// sql: 要执行的sql语句
// ossFilePath: athena会将查询结果保存到一个oss路径上。例如s3://your-bucket/your-prefix
// localPath: 希望保存的文件名，带扩展名。例如xxx.csv
// ossFileDir: 会存在bucketName捅内的 这个路径下
func Download(ak, sk, bucketName, region, sql, ossFilePath, localPath, ossFileDir string) (err error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(ak, sk, ""),
	})
	svc := athena.New(sess)
	input := &athena.StartQueryExecutionInput{
		QueryString: aws.String(sql),
		ResultConfiguration: &athena.ResultConfiguration{
			OutputLocation: aws.String(ossFilePath),
		},
	}

	output, err := svc.StartQueryExecution(input)
	if err != nil {
		return err
	}

	// 下载查询结果
	downloader := s3manager.NewDownloader(sess)
	file, err := os.Create(localPath)
	if err != nil {
		// 创建文件失败
		return err
	}

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),                                           // 替换为你的 S3 桶名
			Key:    aws.String(ossFileDir + "/" + *output.QueryExecutionId + ".csv"), // Athena 会将查询结果保存为 QueryExecutionId.csv
		})

	if err != nil {
		// 下载失败
		return err
	}
	return nil
}

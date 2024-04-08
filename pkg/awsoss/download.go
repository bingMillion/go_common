package awsoss

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"path"
	"strings"
	"time"
)

// Download 根据sql查询athena api，查出的结果保存在目标s3桶中。并下载。
// region 区域，例如 us-west-2
// sql: 要执行的sql语句
// s3ResultDir: s3的结果路径。如"s3://XXXX/testDir"。 XXX为桶名，testDir为桶里的文件名
// localPath: 希望保存的文件名，带扩展名。例如xxx.csv
func Download(ak, sk, region, sql, s3ResultDir, localPath string) (err error) {
	s3Path := strings.TrimPrefix(s3ResultDir, "s3://")
	bucketName, dir := path.Split(s3Path)
	bucketName = strings.TrimSuffix(bucketName, "/")
	dir = strings.TrimSuffix(dir, "/")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(ak, sk, ""),
	})
	svc := athena.New(sess)
	input := &athena.StartQueryExecutionInput{
		QueryString: aws.String(sql),
		ResultConfiguration: &athena.ResultConfiguration{
			OutputLocation: aws.String(s3ResultDir),
		},
	}

	output, err := svc.StartQueryExecution(input) // 只是启动查询，查询仍需要一些时间
	if err != nil {
		return err
	}

	time.Sleep(30 * time.Second)
	output2, err := svc.GetQueryExecution(&athena.GetQueryExecutionInput{
		QueryExecutionId: output.QueryExecutionId,
	})
	if *output2.QueryExecution.Status.State != "SUCCEEDED" {
		return errors.New(*output2.QueryExecution.Status.StateChangeReason)
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
			Bucket: aws.String(bucketName),
			Key:    aws.String(dir + "/" + *output.QueryExecutionId + ".csv"), // 下载路径为目录/QueryExecutionId.csv
		})

	if err != nil {
		// 下载失败
		return err
	}
	return nil
}

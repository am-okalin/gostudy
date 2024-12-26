package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"os"
	"pkg_third/filesystem/s3/conf"
)

func getConf() *conf.AWS {
	return &conf.AWS{
		Region:  "eu-central-1",
		Bucket:  "fstln-scrm-test", //todo::移出该配置
		Workdir: "/test",           //todo::移出该配置
	}
}

func getClient(conf *conf.AWS) *s3.Client {
	//获取s3_config
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(conf.Region),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{Value: aws.Credentials{
			AccessKeyID:     conf.AccessID,
			SecretAccessKey: conf.AccessSecret,
		}}),
	)
	if err != nil {
		log.Fatal(err)
	}

	//获取s3_client
	return s3.NewFromConfig(cfg)
}

//upload 用client将filepath上传到bucket，并命名为key
func upload(client s3.Client, filepath, bucket, key string) error {
	stat, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return err
	}

	//上传文件
	cl := stat.Size()
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: &cl,
	})
	return err
}

func URL(region, bucket, key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
}

//presignURL Get a presigned URL for the object. In order to get a presigned URL for an object, you must create a Presignclient
func presignURL(client s3.Client, bucket, key string) (string, error) {
	//创建预签名对象
	presignClient := s3.NewPresignClient(&client)

	presignResult, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return presignResult.URL, err
}

func download(client s3.Client, bucket, key string) (*s3.GetObjectOutput, error) {
	getObjectResponse, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return getObjectResponse, err
}

func downloadToLocal(client s3.Client, bucket, key string, localPath string) error {
	s3obj, err := download(client, bucket, key)
	if err != nil {
		return err
	}

	file, err := os.Create(localPath)
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(file, s3obj.Body)
	return err
}

func delFilesByBucket(client s3.Client, bucket string) {
	listObjectsV2Response, err := client.ListObjectsV2(context.TODO(),
		&s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
		})
	if err != nil {
		panic("Couldn't list objects...")
	}
	for _, item := range listObjectsV2Response.Contents {
		fmt.Printf("- Deleting object %s\n", *item.Key)
		_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    item.Key,
		})
		if err != nil {
			panic("Couldn't delete items")
		}
	}

	for *listObjectsV2Response.IsTruncated {
		listObjectsV2Response, err = client.ListObjectsV2(context.TODO(),
			&s3.ListObjectsV2Input{
				Bucket:            aws.String(bucket),
				ContinuationToken: listObjectsV2Response.ContinuationToken,
			})
		for _, item := range listObjectsV2Response.Contents {
			fmt.Printf("- Deleting object %s\n", *item.Key)
			_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
				Bucket: aws.String(bucket),
				Key:    item.Key,
			})
			if err != nil {
				panic("Couldn't delete items")
			}
		}
	}
}

func delBucket(client s3.Client, bucket string) {
	_, err := client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		panic("Couldn't delete bucket: " + err.Error())
	}
}

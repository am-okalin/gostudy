package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/xid"
	"os"
	"testing"
)

func TestBucket(t *testing.T) {
	awsConf := getConf()
	s3client := getClient(awsConf)

	myBucketName := "mybucket-" + (xid.New().String())
	MakeBucket(*s3client, myBucketName, awsConf.Region)
	BucketOps(*s3client, myBucketName)
	AccountBucketOps(*s3client, myBucketName)
	BucketDelOps(*s3client, myBucketName)
	//BucketDelOps(*s3client, "mybucket-c9sd3priqgqgho13gqjg")
}

func TestManager(t *testing.T) {
	awsConf := getConf()
	client := getClient(awsConf)

	file, err := os.Open("image.jpg")
	defer file.Close()
	if err != nil {
		panic(err)
	}

	//初始化上传者
	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(awsConf.Bucket),
		Key:    aws.String("my-object-key"),
		Body:   file,
	})
	t.Log(result, err)
}

func TestBucketSign(t *testing.T) {
	conf := getConf()
	client := getClient(conf)

	var err error
	key := "path/myfile.jpg"
	conf.Bucket = "wtn-test"
	//创建bucket
	_, err = client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket:                    aws.String(conf.Bucket),
		ACL:                       types.BucketCannedACLAuthenticatedRead,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{LocationConstraint: types.BucketLocationConstraint(conf.Region)},
	})
	if err != nil {
		panic("could not create bucket: " + err.Error())
	}

	//上传文件
	err = upload(*client, "image.jpg", conf.Bucket, key)
	if err != nil {
		panic("Couldn't upload file: " + err.Error())
	}

	//拼接获取文件
	url1 := URL(conf.Region, conf.Bucket, key)
	t.Log(url1)

	//图片加签
	url2, err := presignURL(*client, conf.Bucket, key)
	if err != nil {
		panic("Couldn't get presigned URL for GetObject")
	}
	t.Log(url2)

	//cleanup
	delFilesByBucket(*client, conf.Bucket)
	delBucket(*client, conf.Bucket)
}

func Test1(t *testing.T) {
	conf := getConf()
	conf.Bucket = "wtn-test"
	client := getClient(conf)
	delFilesByBucket(*client, conf.Bucket)
	delBucket(*client, conf.Bucket)
}

func TestPublicAccessBlock(t *testing.T) {
	conf := getConf()
	conf.Bucket = "wtn-test"
	client := getClient(conf)

	block, err := client.GetPublicAccessBlock(context.TODO(), &s3.GetPublicAccessBlockInput{
		Bucket: aws.String(conf.Bucket),
	})
	t.Log(block, err)
}

func TestBucketPolicy(t *testing.T) {
	conf := getConf()
	client := getClient(conf)

	policy, err := client.GetBucketPolicy(context.TODO(), &s3.GetBucketPolicyInput{
		Bucket: aws.String(conf.Bucket),
	})

	t.Log(policy, err)

	status, err := client.GetBucketPolicyStatus(context.TODO(), &s3.GetBucketPolicyStatusInput{
		Bucket: aws.String(conf.Bucket),
	})

	t.Log(status, err)

}

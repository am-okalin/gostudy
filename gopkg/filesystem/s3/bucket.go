// AWS SDK for Go V2 code examples for Amazon S3
// https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/gov2/s3

package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func MakeBucket(client s3.Client, name string, region string) {
	//snippet-start:[s3.go-v2.CreateBucket]
	// Create a bucket: We're going to create a bucket to hold content.
	// Best practice is to use the preset private access control list (ACL).
	// If you are not creating a bucket from us-east-1, you must specify a bucket location constraint.
	// Bucket names must conform to several rules; read more at https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html
	_, err := client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
		ACL:    types.BucketCannedACLPrivate,
		//ACL:                       types.BucketCannedACLPublicRead,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{LocationConstraint: types.BucketLocationConstraint(region)},
	})

	if err != nil {
		panic("could not create bucket: " + err.Error())
	}

	//snippet-end:[s3.go-v2.CreateBucket]
}

func AccountBucketOps(client s3.Client, name string) {

	fmt.Println("List buckets: ")
	//snippet-start:[s3.go-v2.ListBuckets]
	listBucketsResult, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

	if err != nil {
		panic("Couldn't list buckets")
	}

	for _, bucket := range listBucketsResult.Buckets {
		fmt.Printf("Bucket name: %s\t\tcreated at: %v\n", *bucket.Name, bucket.CreationDate)
	}
	//snippet-end:[s3.go-v2.ListBuckets]

	//snippet-start:[s3.go-v2.ListObjects]
	// List objects in the bucket.
	// n.b. object keys in Amazon S3 do not begin with '/'. You do not need to lead your
	// prefix with it.
	fmt.Println("Listing the objects in the bucket:")
	listObjsResponse, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(name),
		Prefix: aws.String(""),
	})

	if err != nil {
		panic("Couldn't list bucket contents")
	}

	for _, object := range listObjsResponse.Contents {
		fmt.Printf("%s (%d bytes, class %v) \n", *object.Key, object.Size, object.StorageClass)
	}
	//snippet-end:[s3.go-v2.ListObjects]
}

func BucketOps(client s3.Client, bucket string) {
	fmt.Println("Upload an object to the bucket")

	//上传文件
	err := upload(client, "image.jpg", bucket, "path/myfile.jpg")
	if err != nil {
		panic("Couldn't upload file: " + err.Error())
	}

	// 图片加签
	url, err := presignURL(client, bucket, "path/myfile.jpg")
	if err != nil {
		panic("Couldn't get presigned URL for GetObject")
	}
	fmt.Printf("Presigned URL For object: %s\n", url)

	//download
	err = downloadToLocal(client, bucket, "path/myfile.jpg", "download.jpg")
	if err != nil {
		panic(err)
	}

	//copy
	fmt.Println("Copy an object from another bucket to our bucket.")
	_, err = client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(bucket + "/path/myfile.jpg"),
		Key:        aws.String("other/file.jpg"),
	})
	if err != nil {
		panic("Couldn't copy the object to a new key")
	}
}

func BucketDelOps(client s3.Client, name string) {

	//snippet-start:[s3.go-v2.DeleteObject]
	// Delete a single object.
	fmt.Println("Delete an object from a bucket")
	_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(name),
		Key:    aws.String("other/file.jpg"),
	})
	if err != nil {
		panic("Couldn't delete object!")
	}

	//snippet-end:[s3.go-v2.DeleteObject]

	//snippet-start:[s3.go-v2.EmptyBucket]
	// Delete all objects in a bucket.

	fmt.Println("Delete the objects in a bucket")
	// Note: For versioned buckets, you must also delete all versions of
	// all objects within the bucket with ListVersions and DeleteVersion.
	listObjectsV2Response, err := client.ListObjectsV2(context.TODO(),
		&s3.ListObjectsV2Input{
			Bucket: aws.String(name),
		})

	for {

		if err != nil {
			panic("Couldn't list objects...")
		}
		for _, item := range listObjectsV2Response.Contents {
			fmt.Printf("- Deleting object %s\n", *item.Key)
			_, err = client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
				Bucket: aws.String(name),
				Key:    item.Key,
			})

			if err != nil {
				panic("Couldn't delete items")
			}
		}

		if listObjectsV2Response.IsTruncated {
			listObjectsV2Response, err = client.ListObjectsV2(context.TODO(),
				&s3.ListObjectsV2Input{
					Bucket:            aws.String(name),
					ContinuationToken: listObjectsV2Response.ContinuationToken,
				})
		} else {
			break
		}

	}
	//snippet-end:[s3.go-v2.EmptyBucket]

	// snippet-start:[s3.go-v2.DeleteBucket]
	fmt.Println("Delete a bucket")
	// Delete the bucket.

	_, err = client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		panic("Couldn't delete bucket: " + err.Error())
	}
	// snippet-end:[s3.go-v2.DeleteBucket]
}

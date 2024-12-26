package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"strings"
	"testing"
	"time"
)

func TestOss(t *testing.T) {
	// 创建OSSClient实例。
	client, err := oss.New("<yourEndpoint>", "<yourAccessKeyId>", "<yourAccessKeySecret>")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 获取存储空间。
	bucket, err := client.Bucket("<yourBucketName>")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 设置文件元信息：过期时间为2049年1月10日 23:00:00 GMT，访问权限为公共读，自定义元信息为MyProp（取值MyPropVal）。
	expires := time.Date(2049, time.January, 10, 23, 0, 0, 0, time.UTC)
	options := []oss.Option{
		oss.Expires(expires),
		oss.ObjectACL(oss.ACLPublicRead),
		oss.Meta("MyProp", "MyPropVal"),
	}

	// 使用数据流上传文件。
	err = bucket.PutObject("<yourObjectName>", strings.NewReader("MyObjectValue"), options...)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 获取文件元信息。
	props, err := bucket.GetObjectDetailedMeta("<yourObjectName>")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	fmt.Println("Object Meta:", props)
}

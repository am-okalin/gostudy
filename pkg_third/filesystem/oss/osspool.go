package oss

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

type OssPool struct {
	pool     *sync.Pool
	workdir  string
	bucket   string
	endpoint string
}

var lazy *OssPool
var once = &sync.Once{}

//NewOssPool 懒汉单例(lazy load signal)创建ossClient
func NewOssPool(endpoint, bucket, workdir, accessID, accessSecret string) *OssPool {
	once.Do(func() {
		lazy = &OssPool{endpoint: endpoint, bucket: bucket, workdir: workdir, pool: &sync.Pool{New: func() interface{} {
			client, err := oss.New(endpoint, accessID, accessSecret)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			return client
		}}}
	})
	return lazy
}

func (p OssPool) GetWorkDir() string {
	return p.workdir
}

//GetFullPath 获取文件的完整路径
func (p OssPool) GetFullPath(filepath string) string {
	return "https://" + p.bucket + "." + p.endpoint + "/" + p.workdir + filepath
}

//todo::提供bucket供外部使用并回收

//UploadFile 上传文件流
func (p *OssPool) UploadFile(from, to string, options ...oss.Option) error {
	//创建实例
	c := p.pool.Get().(*oss.Client)
	if c == nil {
		return errors.New("初始化oss.client失败")
	}

	//获取存储空间
	bucket, err := c.Bucket(p.bucket)
	if err != nil {
		return err
	}

	//上传文件
	err = bucket.PutObjectFromFile(p.workdir+to, from, options...)
	if err != nil {
		return err
	}

	//将oss连接放回连接池
	p.pool.Put(c)

	return nil
}

//UploadStream 上传文件流
func (p *OssPool) UploadStream(objectName string, r io.Reader) error {
	//创建实例
	c := p.pool.Get().(*oss.Client)
	if c == nil {
		return errors.New("初始化oss.client失败")
	}

	//获取存储空间
	bucket, err := c.Bucket(p.bucket)
	if err != nil {
		return err
	}

	//上传文件
	err = bucket.PutObject(p.workdir+objectName, r)
	if err != nil {
		return err
	}

	//将oss连接放回连接池
	p.pool.Put(c)

	return nil
}

//DownloadStream 下载文件流
func (p *OssPool) DownloadStream(objectName string) (io.ReadCloser, error) {
	//创建实例
	c := p.pool.Get().(*oss.Client)
	if c == nil {
		return nil, errors.New("初始化oss.client失败")
	}

	//获取存储空间
	bucket, err := c.Bucket(p.bucket)
	if err != nil {
		return nil, err
	}

	//下载文件到流
	r, err := bucket.GetObject(p.workdir + objectName)
	if err != nil {
		return nil, err
	}

	//将oss连接放回连接池
	p.pool.Put(c)

	return r, nil
}

//DownloadToFile 下载到本地文件
func (p *OssPool) DownloadToFile(objectName string, localPath string) error {
	//创建实例
	c := p.pool.Get().(*oss.Client)
	if c == nil {
		return errors.New("初始化oss.client失败")
	}

	//获取存储空间
	bucket, err := c.Bucket(p.bucket)
	if err != nil {
		return err
	}

	//下载文件到本地文件中
	err = bucket.GetObjectToFile(p.workdir+objectName, localPath)
	if err != nil {
		return err
	}

	//将oss连接放回连接池
	p.pool.Put(c)

	return nil
}

//GetMate 获取文件信息
func (p *OssPool) GetMate(objectName string) (http.Header, error) {
	//创建实例
	c := p.pool.Get().(*oss.Client)
	if c == nil {
		return nil, errors.New("初始化oss.client失败")
	}

	//获取存储空间
	bucket, err := c.Bucket(p.bucket)
	if err != nil {
		return nil, err
	}

	//下载文件到流
	//mate, err := bucket.GetObjectMeta(p.workdir+file)
	mate, err := bucket.GetObjectDetailedMeta(p.workdir + objectName)
	if err != nil {
		return nil, err
	}

	//将oss连接放回连接池
	p.pool.Put(c)

	return mate, nil
}

type File struct {
	Name         string    `json:"name"`
	FullPath     string    `json:"full_path"`
	Size         int64     `json:"size"`
	SizeStr      string    `json:"size_str"`
	LastModified time.Time `json:"last_modified"`
}

//ListObject 列出目录下所有文件
func (p *OssPool) ListObject(dir string) ([]*File, error) {
	//创建实例
	c := p.pool.Get().(*oss.Client)
	if c == nil {
		return nil, errors.New("初始化oss.client失败")
	}

	//获取存储空间
	bucket, err := c.Bucket(p.bucket)
	if err != nil {
		return nil, err
	}

	//列举所有文件
	var fl []*File
	marker := oss.Marker("")
	prefix := oss.Prefix(p.workdir + dir)
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}
	for {
		lor, err := bucket.ListObjects(marker, prefix)
		if err != nil {
			return nil, err
		}

		for _, object := range lor.Objects {
			fl = append(fl, &File{
				Name:         filepath.Base(object.Key),
				FullPath:     p.GetFullPath(object.Key),
				Size:         object.Size,
				SizeStr:      FormatFileSize(object.Size),
				LastModified: object.LastModified.In(loc),
			})
		}

		prefix = oss.Prefix(lor.Prefix)
		marker = oss.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			break
		}
	}

	//将oss连接放回连接池
	p.pool.Put(c)

	return fl, nil
}

type Files []*File

func (fs Files) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs Files) Len() int {
	return len(fs)
}
func (fs Files) Less(i, j int) bool {
	return fs[i].LastModified.Unix() > fs[j].LastModified.Unix()
}

//FormatFileSize 格式化文件字节大小
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

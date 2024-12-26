package oss

import (
	"io/ioutil"
	"strings"
	"testing"
)

func getConf() (endpoint, bucket, workdir, accessID, accessSecret string) {
	return
}

func Test1(t *testing.T) {
	p := NewOssPool(getConf())
	file := "images/wtn_test.txt"

	// 上传文件
	err := p.UploadStream(file, strings.NewReader("wtn test info1111"))
	if err != nil {
		return
	}

	// 查看元数据
	mate, err := p.GetMate(file)
	if err != nil {
		return
	}
	t.Log(mate)

	// 下载文件
	stream, err := p.DownloadStream(file)
	if err != nil {
		return
	}
	defer stream.Close()
	// 解析下载文件
	var res []byte
	res, err = ioutil.ReadAll(stream)
	if err != nil {
		return
	}
	t.Log(string(res))
}

func Test2(t *testing.T) {
	p := NewOssPool(getConf())
	dir := "webmap/v"
	_, err := p.ListObject(dir)
	if err != nil {
		return
	}
}

func Test_ossPool_DownloadToFile(t *testing.T) {
	url := "webmap/v0.1.0.tar.gz"
	localPath := "C:/Users/tony/Downloads/b.tar.gz"
	p := NewOssPool(getConf())
	err := p.DownloadToFile(url, localPath)
	if err != nil {
		t.Error(err)
	}
}

func Test_ossPool_UploadFile(t *testing.T) {
	localPath := "C:/Users/tony/Downloads/b.tar.gz"
	p := NewOssPool(getConf())
	err := p.UploadFile(localPath, "webmap/b.tar.gz")
	if err != nil {
		t.Error(err)
	}
}

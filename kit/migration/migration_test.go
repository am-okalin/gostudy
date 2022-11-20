package migration

import (
	"testing"
	"time"
)

const (
	fname = "./file/worker.csv"
)

type Worker struct {
	Id          int64
	Name        string
	Description string
	CreateTime  time.Time
	DeleteTime  time.Time
}

func Test1(t *testing.T) {
	var list []Worker
	err := Unmarshal(fname, &list)
	if err != nil {
		return
	}

	newFile, err := Marshal(list)
	if err != nil {
		return
	}
	t.Log(newFile)
}

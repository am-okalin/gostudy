package option

import (
	"fmt"
	"testing"
)

type StuffClient interface {
	DoStuff() error
}

type stuffClient struct {
	conn Connection
	opts *Options
}

type Connection struct{}

func NewStuffClient(conn Connection, options ...Option) StuffClient {
	return &stuffClient{
		conn: conn,
		opts: loadOptions(options...),
	}
}

func (c stuffClient) DoStuff() error {
	fmt.Println("client do stuff")
	return nil
}

func TestNewStuffClient(t *testing.T) {
	client := NewStuffClient(
		Connection{},
		WithTimeout(10),
		WithRetries(10),
	)
	err := client.DoStuff()
	if err != nil {
		t.Error(err)
	}
}

package sync

import (
	"bytes"
	"fmt"
	"github.com/alitto/pond"
	"sync"
	"time"
)

var buffers = sync.Pool{New: func() interface{} {
	return new(bytes.Buffer)
}}

func GetBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}
func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	buffers.Put(buf)
}

func wp() {

	// Create a buffered (non-blocking) pool that can scale up to 100 workers
	// and has a buffer capacity of 1000 tasks
	pool := pond.New(10, 100)

	// Submit 1000 tasks
	for i := 0; i < 100; i++ {
		n := i
		pool.Submit(func() {
			time.Sleep(1 * time.Second)
			fmt.Printf("Running task #%d\n", n)
		})
	}

	// Stop the pool and wait for all submitted tasks to complete
	pool.StopAndWait()
}

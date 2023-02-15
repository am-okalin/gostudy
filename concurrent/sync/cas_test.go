package sync

import (
	"fmt"
	"sync/atomic"
	"testing"
	"unsafe"
)

var global [1024]byte

func testPointers() []unsafe.Pointer {
	var pointers []unsafe.Pointer
	// globals
	for i := 0; i < 10; i++ {
		pointers = append(pointers, unsafe.Pointer(&global[1<<i-1]))
	}
	// heap
	pointers = append(pointers, unsafe.Pointer(new(byte)))
	// nil
	pointers = append(pointers, nil)
	return pointers
}

func TestSwapPointer(t *testing.T) {
	var x struct {
		before uintptr
		i      unsafe.Pointer
		after  uintptr
	}
	var m uint64 = 1
	magicptr := uintptr(m)
	x.before = magicptr
	x.after = magicptr
	var j unsafe.Pointer

	for _, p := range testPointers() {
		k := atomic.SwapPointer(&x.i, p)
		if x.i != p || k != j {
			t.Fatalf("p=%p i=%p j=%p k=%p", p, x.i, j, k)
		}
		j = p
	}
	if x.before != magicptr || x.after != magicptr {
		t.Fatalf("wrong magic: %#x _ %#x != %#x _ %#x", x.before, x.after, magicptr, magicptr)
	}
}

func TestSwapPointerM(t *testing.T) {
	newM := map[string]uint64{}
	nowM := map[string]uint64{"kk": 1}
	nowP := unsafe.Pointer(&nowM)
	oldP := atomic.SwapPointer(&nowP, unsafe.Pointer(&newM))

	now := *(*map[string]uint64)(nowP)
	old := *(*map[string]uint64)(oldP)
	fmt.Printf("map: %v\n", now)
	fmt.Printf("map: %v\n", old)
}

func SwarpM(nowM, newM map[string]uint64) map[string]uint64 {
	nowP := unsafe.Pointer(&nowM)
	oldP := atomic.SwapPointer(&nowP, unsafe.Pointer(&newM))
	return *(*map[string]uint64)(oldP)
}

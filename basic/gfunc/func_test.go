package gfunc

import (
	"runtime"
	"testing"
)

func Test1(t *testing.T) {
	f1()
}

func Test2(t *testing.T) {
	f1 := getEatFruitFun()
	w1 := f1("apple")
	t.Log(w1)

	f2 := getEatMeatFun()
	w2 := f2("Beef")
	t.Log(w2)
}

var oneMB []byte

func TestCountMallocs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping malloc count in short mode")
	}
	if runtime.GOMAXPROCS(0) > 1 {
		t.Skip("skipping; GOMAXPROCS>1")
	}
	// Allocate a big messy buffer for AppendQuoteToASCII's test.
	oneMB = make([]byte, 1e6)
	for i := range oneMB {
		oneMB[i] = byte(i)
	}
	for _, mt := range mallocTest {
		allocs := testing.AllocsPerRun(100, mt.fn)
		if max := float64(mt.count); allocs > max {
			t.Errorf("%s: %v allocs, want <=%v", mt.desc, allocs, max)
		}
	}
}

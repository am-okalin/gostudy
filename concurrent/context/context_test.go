package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestValueCtx(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, "k", "v1")
	ctx = context.WithValue(ctx, "k", "v2")
	t.Log(ctx.Value("k"))
}

type otherContext struct {
	context.Context
}

func TestWithCancel(t *testing.T) {
	c0 := context.Background()
	c1, cancel1 := context.WithCancel(c0)
	defer cancel1()
	c2, cancel2 := context.WithCancel(c1)
	defer cancel2()

	if got, want := fmt.Sprint(c2), "context.Background.WithCancel.WithCancel"; got != want {
		t.Errorf("c2.String() = %q want %q", got, want)
	}

	//即便链路中存在非cancelCtx的类型，祖先CancelCtx被Cancel时也会通知到子孙CancelCtx
	c3 := otherContext{c2}
	c4, cancel4 := context.WithCancel(c3)
	defer cancel4()

	// c0 是 Background，Done 为 nil；c1~c4 都有各自的 Done
	if d := c0.Done(); d != nil {
		t.Errorf("c0.Done() == %v want nil", d)
	}
	for i, c := range []context.Context{c1, c2, c3, c4} {
		if d := c.Done(); d == nil {
			t.Errorf("c[%d].Done() == %v want non-nil", i+1, d)
		}
		if e := c.Err(); e != nil {
			t.Errorf("c[%d].Err() == %v want nil", i+1, e)
		}

		select {
		case x := <-c.Done():
			t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
		default:
		}
	}

	cancel2() // Should propagate synchronously.
	// cancel2 会同步通知到 c2 及其子孙 c3、c4，但不会影响祖先 c0、c1
	for i, c := range []context.Context{c2, c3, c4} {
		select {
		case <-c.Done():
		default:
			t.Errorf("<-c[%d].Done() blocked, but shouldn't have", i+2)
		}
		if e := c.Err(); e != context.Canceled {
			t.Errorf("c[%d].Err() == %v want %v", i+2, e, context.Canceled)
		}
	}
	for i, c := range []context.Context{c0, c1} {
		if e := c.Err(); e != nil {
			t.Errorf("c[%d].Err() == %v want nil", i, e)
		}
	}
}

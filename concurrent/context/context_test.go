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
	c1, _ := context.WithCancel(c0)
	c2, cancel2 := context.WithCancel(c1)

	if got, want := fmt.Sprint(c2), "context.Background.WithCancel.WithCancel"; got != want {
		t.Errorf("c1.String() = %q want %q", got, want)
	}

	//即便链路中存在非cancelCtx的类型，祖先CancelCtx被Cancel时也会通知到子孙CancelCtx
	c3 := otherContext{c2}
	c4, _ := context.WithCancel(c3)
	contexts := []context.Context{c0, c1, c2, c3, c4}

	for i, c := range contexts {
		if d := c.Done(); d == nil {
			//c0.Done返回nil, c3.Done则返回c2.Done对象, 其他返回各自的Done对象
			t.Errorf("c[%d].Done() == %v want non-nil", i, d)
		}
		if e := c.Err(); e != nil {
			t.Errorf("c[%d].Err() == %v want nil", i, e)
		}

		select {
		case x := <-c.Done():
			t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
		default:
		}
	}

	cancel2() // Should propagate synchronously.
	for i, c := range contexts {
		select {
		case <-c.Done():
		default:
			//cancel2无法影响祖先，所以c0,c1仍然blocked
			t.Errorf("<-c[%d].Done() blocked, but shouldn't have", i)
		}
		if e := c.Err(); e != context.Canceled {
			t.Errorf("c[%d].Err() == %v want %v", i, e, context.Canceled)
		}
	}
}

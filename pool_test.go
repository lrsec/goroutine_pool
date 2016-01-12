package goroutine_pool

import (
	"fmt"
	"testing"
	"time"
)

func TestNewWithError(t *testing.T) {
	_, err := NewPool(-1, 1, 1, 1, 1, 1, func(interface{}) interface{} {
		return nil
	})

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}

	_, err = NewPool(1, -1, 1, 1, 1, 1, func(interface{}) interface{} {
		return nil
	})

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}

	_, err = NewPool(1, 1, -1, 1, 1, 1, func(interface{}) interface{} {
		return nil
	})

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}

	_, err = NewPool(1, 1, 1, -1, 1, 1, func(interface{}) interface{} {
		return nil
	})

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}

	_, err = NewPool(1, 1, 1, 1, -1, 1, func(interface{}) interface{} {
		return nil
	})

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}

	_, err = NewPool(1, 1, 1, 1, 1, -1, func(interface{}) interface{} {
		return nil
	})

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}

	_, err = NewPool(1, 1, 1, 1, 1, 1, nil)

	fmt.Println(err)
	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}
}

func TestNewPool(t *testing.T) {

	pool, err := NewPool(10, 20, 15, 500, 10, 10, func(input interface{}) interface{} {
		return input
	})

	if err != nil {
		t.Error("Can not create pool", err)
	}

	defer func() {
		if pool != nil {
			pool.Close()
		}
	}()

	for i := 0; i < 10; i++ {
		pool.InboundChannel() <- i
	}

	time.Sleep(2 * time.Second)

	length := len(pool.OutboundChannel())

	if length != 10 {
		t.Errorf("data not treated correctly with outbound len: %d\n", length)
	}
}

package goroutine_pool

import (
	"testing"
	"time"
)

func TestNewWithError(t *testing.T) {
	inboundChannel := make(chan interface{})

	defer close(inboundChannel)

	_, err := NewPool(-1, 1, 1, 1, 1, 1, inboundChannel, func(interface{}) interface{} {
		return nil
	})

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}

	_, err = NewPool(1, -1, 1, 1, 1, 1, inboundChannel, func(interface{}) interface{} {
		return nil
	})

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}

	_, err = NewPool(1, 1, -1, 1, 1, 1, inboundChannel, func(interface{}) interface{} {
		return nil
	})

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}

	_, err = NewPool(1, 1, 1, -1, 1, 1, inboundChannel, func(interface{}) interface{} {
		return nil
	})

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}

	_, err = NewPool(1, 1, 1, 1, -1, 1, inboundChannel, func(interface{}) interface{} {
		return nil
	})

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}

	_, err = NewPool(1, 1, 1, 1, 1, -1, inboundChannel, func(interface{}) interface{} {
		return nil
	})

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}

	_, err = NewPool(1, 1, 1, 1, 1, 1, nil, func(interface{}) interface{} {
		return nil
	})

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}

	_, err = NewPool(1, 1, 1, 1, 1, 1, inboundChannel, nil)

	if err == nil {
		t.Error("NewPool does not report error for incorrect params")
	}
}

func TestNewPool(t *testing.T) {

	inboundChannel := make(chan interface{}, 10)
	defer close(inboundChannel)

	pool, err := NewPool(10, 20, 15, 500, 10, 10, inboundChannel, func(input interface{}) interface{} {
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
		pool.InboundChannel <- i
	}

	time.Sleep(2 * time.Second)

	length := len(pool.OutboundChannel)

	if length != 10 {
		t.Errorf("data not treated correctly with outbound len: %d\n", length)
	}
}

func TestWithPanic(t *testing.T) {
	inboundChannel := make(chan interface{}, 10)
	defer close(inboundChannel)

	pool, err := NewPool(10, 20, 15, 500, 10, 10, inboundChannel, func(input interface{}) interface{} {
		panic("test")
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
		pool.InboundChannel <- i
	}

	time.Sleep(2 * time.Second)

	inLength := len(pool.InboundChannel)
	outLength := len(pool.OutboundChannel)

	if inLength != 0 {
		t.Errorf("data not treated correctly with inbound len: %d\n", inLength)
	}

	if outLength != 0 {
		t.Errorf("data not treated correctly with outbound len: %d\n", outLength)
	}
}

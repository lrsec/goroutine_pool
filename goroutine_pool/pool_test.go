package goroutine_pool

import (
	"testing"
	"time"
)

func TestNewWithError(t *testing.T) {
	inboundChannel := make(chan interface{}, 10)
	outboundChannel := make(chan interface{}, 10)

	defer func() {
		close(inboundChannel)
		close(outboundChannel)
	}()

	_, err := NewGPool(-1, 1, 1, 1*time.Second, 1*time.Second, inboundChannel, outboundChannel, func(interface{}) (interface{}, error) {
		return nil, nil
	}, "test")

	if err == nil {
		t.Error("NewGPool does not report error for incorrect params")
	}

	_, err = NewGPool(1, -1, 1, 1*time.Second, 1*time.Second, inboundChannel, outboundChannel, func(interface{}) (interface{}, error) {
		return nil, nil
	}, "test")

	if err == nil {
		t.Error("NewGPool does not report error for incorrect params")
	}

	_, err = NewGPool(1, 1, -1, 1*time.Second, 1*time.Second, inboundChannel, outboundChannel, func(interface{}) (interface{}, error) {
		return nil, nil
	}, "test")

	if err == nil {
		t.Error("NewGPool does not report error for incorrect params")
	}

	_, err = NewGPool(1, 1, 1, 1*time.Second, 1*time.Second, nil, outboundChannel, func(interface{}) (interface{}, error) {
		return nil, nil
	}, "test")

	if err == nil {
		t.Error("NewGPool does not report error for incorrect params")
	}

	_, err = NewGPool(1, 1, 1, 1*time.Second, 1*time.Second, inboundChannel, nil, func(interface{}) (interface{}, error) {
		return nil, nil
	}, "test")

	if err == nil {
		t.Error("NewGPool does not report error for incorrect params")
	}

	_, err = NewGPool(1, 1, 1, 1*time.Second, 1*time.Second, inboundChannel, outboundChannel, nil, "test")

	if err == nil {
		t.Error("NewGPool does not report error for incorrect params")
	}
}

func TestNewGPool(t *testing.T) {

	inboundChannel := make(chan interface{}, 10)
	outboundChannel := make(chan interface{}, 10)

	defer func() {
		close(inboundChannel)
		close(outboundChannel)
	}()

	pool, err := NewGPool(10, 20, 15, 500*time.Millisecond, 500*time.Millisecond, inboundChannel, outboundChannel, func(input interface{}) (interface{}, error) {
		return input, nil
	}, "test")

	if err != nil {
		t.Error("Can not create pool", err)
	}

	defer func() {
		if pool != nil {
			pool.Close()
		}
	}()

	for i := 0; i < 10; i++ {
		pool.InputChannel <- i
	}

	time.Sleep(2 * time.Second)

	length := len(pool.OutputChannel)

	if length != 10 {
		t.Errorf("data not treated correctly with outbound len: %d\n", length)
	}
}

func TestWithPanic(t *testing.T) {
	inboundChannel := make(chan interface{}, 10)
	outboundChannel := make(chan interface{}, 10)

	defer func() {
		close(inboundChannel)
		close(outboundChannel)
	}()

	pool, err := NewGPool(10, 20, 15, 500*time.Millisecond, 500*time.Millisecond, inboundChannel, outboundChannel, func(input interface{}) (interface{}, error) {
		panic("test")
	}, "test")

	if err != nil {
		t.Error("Can not create pool", err)
	}

	defer func() {
		if pool != nil {
			pool.Close()
		}
	}()

	for i := 0; i < 10; i++ {
		pool.InputChannel <- i
	}

	time.Sleep(2 * time.Second)

	inLength := len(pool.InputChannel)
	outLength := len(pool.OutputChannel)

	if inLength != 0 {
		t.Errorf("data not treated correctly with inbound len: %d\n", inLength)
	}

	if outLength != 0 {
		t.Errorf("data not treated correctly with outbound len: %d\n", outLength)
	}
}

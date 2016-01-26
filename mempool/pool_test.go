package mempool

import (
	"fmt"
	"github.com/lrsec/pool/goroutine_pool"
	"runtime"
	"testing"
	"time"
)

func TestLargerPowerOf2(t *testing.T) {

	fmt.Println(runtime.GOMAXPROCS(0))

	p := LargerPowerOf2(0)
	if p != 0 {
		t.Fatalf("largerPowerOf2 does not fit for 0. expacted: 0, got: %d", p)
	}

	p = LargerPowerOf2(1)
	if p != 0 {
		t.Fatalf("largerPowerOf2 does not fit for 0. expacted: 0, got: %d", p)
	}

	var i uint64

	inbound := make(chan interface{}, 10*1000*1000)

	var pool *goroutine_pool.Pool
	pool, err := goroutine_pool.NewPool(100000, 100000, 150000, 5*60*1000, 100, 1000000, inbound, func(msg interface{}, outbound chan<- interface{}) interface{} {
		message := msg.(Message)

		start := time.Now().UnixNano()
		power := LargerPowerOf2(message.value)
		end := time.Now().UnixNano()

		fmt.Printf("Dealing value %d, power: %d, cost time: %dms \n", message.value, message.power, end-start)

		if uint64(power) != message.power {
			fmt.Errorf("largerPowerOf2 does not fit for %d. expacted: %d, got: %d", message.value, message.power, power)
			outbound <- false
			pool.Close()
		}

		return nil
	})

	if err != nil {
		t.Fatal("Can not produce goroutine pool", err.Error())
	}

	for i = 1; i < 16; i++ {

		var max uint64 = 1 << i
		var min uint64 = 1<<(i-1) + 1

		for j := min; j <= max; j++ {

			msg := Message{
				value: j,
				power: i,
			}

			select {
			case inbound <- msg:
			case <-pool.PoolCloseSignal:
				t.Fail()
			}
		}
	}

	t.Log("Submition Complete")

	select {
	case _, ok := <-pool.OutboundChannel:
		if ok {
			t.Fail()
		} else {
			if len(pool.OutboundChannel) > 0 {
				t.Fail()
			}
		}
	case <-pool.PoolCloseSignal:
		if len(pool.OutboundChannel) > 0 {
			t.Fail()
		}
	case <-time.After(1 * 60 * time.Second):
		t.Log("Test complete")

	}

}

type Message struct {
	value uint64
	power uint64
}

package goroutine_pool

import (
	"fmt"
	log "github.com/cihub/seelog"
	errors "github.com/go-errors/errors"
	"sync/atomic"
	"time"
)

func NewPool(initPoolSize, maxPoolSize, maxIdleSize, maxIdleMs, inboundChannelBuffer, outboundChannelBuffer int64, handler func(interface{}) interface{}) (*Pool, error) {
	if initPoolSize < 0 || maxPoolSize < 0 || maxIdleSize < 0 || maxIdleMs < 0 || inboundChannelBuffer < 0 || outboundChannelBuffer < 0 || handler == nil {
		return nil, errors.New(fmt.Sprintf("Illegal parameters to create goroutine pool. initPoolSize: %v, maxPoolSize: %v , maxIdleSize: %v, maxIdleMs: %v, inboundChannelBuffer: %v, outboundChannelBuffer: %v, handler: %v", initPoolSize, maxPoolSize, maxIdleSize, maxIdleMs, inboundChannelBuffer, outboundChannelBuffer, handler))
	}

	pool := &Pool{
		defaultInitPoolSize: initPoolSize,
		defaultMaxPoolSize:  maxPoolSize,
		defaultMaxIdleSize:  maxIdleSize,
		defaultMaxIdleMs:    maxIdleMs,

		defaultMonitorMs: 1000,

		inboundChannel:  make(chan interface{}, inboundChannelBuffer),
		outboundChannel: make(chan interface{}, outboundChannelBuffer),

		poolCloseSignal: make(chan bool),
	}

	pool.start(handler)

	return pool, nil
}

type Pool struct {
	defaultInitPoolSize int64
	defaultMaxPoolSize  int64
	defaultMaxIdleMs    int64
	defaultMaxIdleSize  int64

	defaultMonitorMs int64

	poolSize int64

	inboundChannel  chan interface{}
	outboundChannel chan interface{}

	poolCloseSignal chan bool
}

func (pool *Pool) start(handler func(interface{}) interface{}) error {

	// worker definition
	worker := func() {

		atomic.AddInt64(&pool.poolSize, 1)

		for {
			select {
			case c := <-pool.inboundChannel:
				result := func() interface{} {
					defer func() {
						if r := recover(); r != nil {
							log.Error("woker handler panic", r)
							return nil
						}
					}()

					return handler(c)
				}()
				if result != nil {
					pool.outboundChannel <- result
				}
			case <-pool.poolCloseSignal:
				log.Trace("Recive close signal. Close go routine.")
				atomic.AddInt64(&pool.poolSize, -1)
				return
			case <-time.After(pool.defaultMaxIdleMs * time.Millisecond):
				size := atomic.LoadInt64(pool.poolSize)
				if size > pool.defaultMaxIdleSize && atomic.CompareAndSwapInt64(pool.poolSize, size, size-1) {
					log.Trace("Pool idle time longer than maxIdleTime, close")
					return
				}

			}
		}
	}

	// start workers
	for i := 0; i < pool.defaultInitPoolSize; i++ {
		go worker()
	}

	// start supervisor
	go func() {
		for {
			select {
			case <-pool.poolCloseSignal:
				log.Debug("Goroutine is closed. Close all goroutines")
				return
			case <-time.After(pool.defaultMonitorMs * time.Millisecond):

				blockedSize := len(pool.inboundChannel)
				if blockedSize > 0 {
					poolSize := atomic.LoadInt64(pool.poolSize)
					if poolSize < pool.defaultMaxPoolSize {
						size := poolSize - pool.defaultMaxPoolSize
						if size > blockedSize {
							size = blockedSize
						}

						for i := 0; i < size; i++ {
							go worker()
						}
					}
				}
			}
		}
	}()
}

func (pool *Pool) close() {
	close(pool.poolCloseSignal)
}

func (pool *Pool) Produce(task interface{}) {
	pool.inboundChannel <- task
}

func (pool *Pool) Consume() (interface{}, bool) {
	return <-pool.outboundChannel
}

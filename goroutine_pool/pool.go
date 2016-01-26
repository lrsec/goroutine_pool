package goroutine_pool

import (
	"fmt"
	log "github.com/cihub/seelog"
	errors "github.com/lrsec/errors/wrapper"
	"sync/atomic"
	"time"
)

func NewPool(initPoolSize, maxPoolSize, maxIdleSize, maxIdleMs, outboundChannelBuffer, watermark int64, inboundChannel chan interface{}, handler func(interface{}, chan<- interface{}) interface{}) (*Pool, error) {
	if initPoolSize < 0 || maxPoolSize < 0 || maxIdleSize < 0 || maxIdleMs < 0 || outboundChannelBuffer < 0 || watermark < 0 || inboundChannel == nil || handler == nil {
		return nil, errors.New(fmt.Sprintf("Illegal parameters to create goroutine pool. initPoolSize: %v, maxPoolSize: %v , maxIdleSize: %v, maxIdleMs: %v, outboundChannelBuffer: %v, watermark: %v, inboundChannel: %v handler: %v", initPoolSize, maxPoolSize, maxIdleSize, maxIdleMs, outboundChannelBuffer, watermark, inboundChannel, handler))
	}

	pool := &Pool{
		initPoolSize: initPoolSize,
		maxPoolSize:  maxPoolSize,
		maxIdleSize:  maxIdleSize,
		maxIdleMs:    maxIdleMs,
		watermark:    watermark,

		monitorMs: 1000,

		InboundChannel:  inboundChannel,
		OutboundChannel: make(chan interface{}, outboundChannelBuffer),

		PoolCloseSignal: make(chan bool),
	}

	pool.start(handler)

	return pool, nil
}

type Pool struct {
	initPoolSize int64
	maxPoolSize  int64
	maxIdleMs    int64
	maxIdleSize  int64
	monitorMs    int64
	watermark    int64

	poolSize int64

	InboundChannel  chan interface{}
	OutboundChannel chan interface{}

	PoolCloseSignal chan bool
}

func (pool *Pool) start(handler func(interface{}, chan<- interface{}) interface{}) {

	// worker definition
	worker := func() {
		atomic.AddInt64(&pool.poolSize, 1)
		timer := time.NewTimer(0)

		for {
			timer.Reset(time.Duration(pool.maxIdleMs) * time.Millisecond)

			select {
			case c := <-pool.InboundChannel:
				result := func() interface{} {
					defer func() {
						if r := recover(); r != nil {
							log.Error("woker handler panic", r)
						}
					}()

					return handler(c, pool.OutboundChannel)
				}()
				if result != nil {
					pool.OutboundChannel <- result
				}
			case <-pool.PoolCloseSignal:
				log.Trace("Recive close signal. Close go routine.")
				atomic.AddInt64(&pool.poolSize, -1)
				return
			case <-timer.C:
				size := atomic.LoadInt64(&pool.poolSize)
				if size > pool.maxIdleSize && atomic.CompareAndSwapInt64(&pool.poolSize, size, size-1) {
					log.Trace("Pool idle time longer than maxIdleTime, close")
					return
				}

			}
		}
	}

	// start workers
	for i := int64(0); i < pool.initPoolSize; i++ {
		go worker()
	}

	// start supervisor
	go func() {
		timer := time.NewTimer(0)

		for {
			timer.Reset(time.Duration(pool.monitorMs) * time.Millisecond)

			select {
			case <-pool.PoolCloseSignal:
				log.Debug("Goroutine is closed. Close all goroutines")
				close(pool.OutboundChannel)
				return
			case <-timer.C:
				blockedSize := int64(len(pool.InboundChannel))
				if blockedSize > pool.watermark {
					poolSize := atomic.LoadInt64(&pool.poolSize)
					if poolSize < pool.maxPoolSize {
						size := poolSize - pool.maxPoolSize
						blockedSize *= 2

						if size > blockedSize {
							size = blockedSize
						}

						for i := int64(0); i < size; i++ {
							go worker()
						}
					}
				}
			}
		}
	}()
}

func (pool *Pool) Close() {
	close(pool.PoolCloseSignal)
}

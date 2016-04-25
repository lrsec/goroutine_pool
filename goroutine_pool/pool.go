package goroutine_pool

import (
	"fmt"
	log "github.com/cihub/seelog"
	errors "github.com/lrsec/errors/wrapper"
	"sync/atomic"
	"time"
)

type GPool struct {
	minSize     int64
	maxSize     int64
	maxIdleSize int64
	maxIdleTime time.Duration

	watermark int64

	monitorPeriod time.Duration

	poolSize int64
	isClosed atomic.Value

	InputChannel  chan interface{}
	OutputChannel chan interface{}

	Name string
}

func NewGPool(minSize, maxSize, maxIdleSize int64,
	maxIdleTime, monitorPeriod time.Duration,
	inputChannel, outputChannel chan interface{},
	handler func(interface{}) (interface{}, error),
	name string) (*GPool, error) {

	if minSize < 0 || maxSize < 0 || maxIdleSize < 0 || handler == nil || inputChannel == nil || outputChannel == nil {
		return nil, errors.New(fmt.Sprintf("Illegal parameters to create goroutine pool. minSize: %v, maxSize: %v, maxIdleSize: %v, handler: %v. inputChannel: %v. outputChannel: %v", minSize, maxSize, maxIdleSize, handler, inputChannel, outputChannel))
	}

	gpool := &GPool{
		minSize:     minSize,
		maxSize:     maxSize,
		maxIdleSize: maxIdleSize,
		maxIdleTime: maxIdleTime,

		monitorPeriod: monitorPeriod,

		poolSize:  0,
		watermark: minSize / 2,

		InputChannel:  inputChannel,
		OutputChannel: outputChannel,

		Name: name,
	}

	gpool.isClosed.Store(false)

	gpool.start(handler)

	return gpool, nil
}

func (gpool *GPool) start(handler func(interface{}) (interface{}, error)) {

	// worker definition
	worker := func() {
		atomic.AddInt64(&(gpool.poolSize), 1)
		defer atomic.AddInt64(&(gpool.poolSize), -1)

		timer := time.NewTimer(0)

		for {
			timer.Reset(gpool.maxIdleTime)

			select {
			case c := <-gpool.InputChannel:
				result, err := func() (interface{}, error) {
					defer func() {
						if r := recover(); r != nil {
							log.Error("woker handler panic", r)
						}
					}()

					return handler(c)
				}()
				if err != nil {
					log.Warnf("GPool %s handler return error for input: %v. Error: %s", gpool.Name, err.Error())
				} else {
					if result != nil {
						gpool.OutputChannel <- result
					} else {
						log.Warnf("GPool %s handler get an nil result for input: %v", gpool.Name, c)
					}
				}

			case <-timer.C:
				poolSize := atomic.LoadInt64(&(gpool.poolSize))

				if poolSize > gpool.maxIdleSize && atomic.CompareAndSwapInt64(&(gpool.poolSize), poolSize, poolSize-1) {
					log.Trace("Pool idle time longer than maxIdleTime, close")
					return
				}

			}
		}
	}

	// start workers
	for i := int64(0); i < gpool.minSize; i++ {
		go worker()
	}

	// start supervisor
	go func() {
		timer := time.NewTimer(0)

		for {
			// 如果池关闭,则不在进行 goroutines 扩展与调度
			if gpool.isClosed.Load().(bool) {
				return
			}

			timer.Reset(gpool.monitorPeriod)

			select {

			case <-timer.C:
				blockedSize := int64(len(gpool.InputChannel))

				if blockedSize > gpool.watermark {
					poolSize := atomic.LoadInt64(&(gpool.poolSize))
					if poolSize < gpool.maxSize {
						size := gpool.maxSize - poolSize
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

func (gpool *GPool) Close() {
	gpool.isClosed.Store(true)
}

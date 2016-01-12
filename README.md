# goroutine_pool

goroutine_pool is a dynamic pool of goroutines followed the pipeline pattern.

# Purpose

Long-live and reusable goroutines are created in pool, instead of short-live goroutines, in order to reduce the overheads of goroutines and the burden of GC.

# Usage

get the package

```
go get -u github.com/lrsec/goroutine_pool
```

Get the pool

```go

import "github.com/lrsec/goroutine_pool"

pool, err := goroutine_pool.NewPool(1000, 10000, 1500, 10 * 60 * 1000, 100, 100, func (input interface{}) interface{} {
    return input
})


```

Use the inbound and outbound channels

```
inbound := pool.InboundChannel()

outbound := pool.OutboundChannel()
```


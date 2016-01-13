# goroutine_pool

goroutine_pool is a dynamic pool of goroutines followed the pipeline pattern.

# Purpose
* Reuse goroutines
* Reduce overheads of creating huge number of short-live goroutines
* Reduce burden of GC
* Control total number of goroutines

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


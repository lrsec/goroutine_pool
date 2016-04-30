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

inputChannel chan interface{}
outputChannel chan interface{}

pool, err := goroutine_pool.NewPool(100, 10000, 150, 60 * time.Second, 5 * time.Second, inputChannel, outputChannel, "test_pool", func (message interface{}, channel chan interface{}) error {
    outputChannel<-message
    return nil
})


```



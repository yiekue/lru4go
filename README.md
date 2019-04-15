# lru4go
LRU cache implementation based on go, supports expiration time.

+ thread safety
+ support expiration time
+ lazy elimination

example:
```go
package main

import "lru4go"

func main() {
    cache := lru4go.New(2000)
    cache.Set("key","value")
    cache.Set("ttl","20 second", 20)
    if v, ok:=cache.Get(key); !ok {
        fmt.Println("get failed:", ok.Error())
    }
    _ := cache.Delete("ttl")
    cache.Reset()
    cache.Set("key","value")
}
```
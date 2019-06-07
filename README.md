# lru4go
使用golang实现的LRU缓存.

+ 线程安全
+ 支持设置过期时间

example:
```go

cache := lru4go.New(2000)
cache.Set("key","value")
cache.Set("ttl","20 second", 20)
if v, ok:=cache.Get(key); !ok {
    fmt.Println("get failed:", ok.Error())
}
_ := cache.Delete("ttl")
cache.Reset()
cache.Set("key","value")

```
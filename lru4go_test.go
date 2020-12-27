package lru4go

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	cache, err := New(50)
	if err != nil || cache == nil {
		t.Fatal("create cache failed!")
	}
}

func TestLrucache_Set(t *testing.T) {
	cache, _ := New(50)
	if cache == nil {
		t.Fatal("create cache failed!")
	}
	cache.Set("test1", 123)
	v, err := cache.Get("test1")
	if err != nil {
		t.Fatal("get failed,err:", err)
	}
	if v == nil {
		t.Fatal("get failed, value is nil")
	}
}

func TestLrucache_Keys(t *testing.T) {
	cache, _ := New(50)
	cache.Set("t1", 1)
	cache.Set("t2", 2)
	cache.Set("t3", 2)

	keys := cache.Keys()
	if nil == keys {
		t.Fatal("keys is nil")
	}
	for _, v := range keys {
		if "t1" != v && "t2" != v && "t3" != v {
			t.Error("keys wrong,key=", v)
		}
	}
}

func TestLrucache_Delete(t *testing.T) {
	cache, _ := New(50)
	cache.Set("t1", 1)
	cache.Set("t2", 2)

	cache.Delete("t1")

	v, _ := cache.Get("t1")
	if v != nil {
		t.Error("delete failed!")
	}
}

func TestLrucache_expireOldest(t *testing.T) {
	max := 10
	cache, _ := New(max)
	for i := 0; i <= max+1; i++ {
		cache.Set(fmt.Sprintf("%d", i), i)
	}

	v, err := cache.Get(fmt.Sprintf("%d", 0))
	if err == nil {
		t.Error("expire oldest failed:", v)
	}
}

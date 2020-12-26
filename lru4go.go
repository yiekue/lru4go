package lru4go

import (
	"container/list"
	"errors"
	"sync"
)

// Lrucache cache data structure
type Lrucache struct {
	cap  int
	dict map[string]*list.Element
	l    *list.List
	mu   *sync.Mutex
}

type elem struct {
	k      string
	v      interface{}
	expire int64
}

// New create a new lrucache
// size: max number of element
func New(size int) (*Lrucache, error) {
	if size < 0 {
		return nil, errors.New("size must be positive")
	}
	lc := &Lrucache{
		cap:  size,
		l:    list.New(),
		dict: make(map[string]*list.Element),
		mu:   &sync.Mutex{},
	}
	return lc, nil
}

// Set create or update an element using key
// 		key:	The identity of an element
// 		value: 	new value of the element
func (lc *Lrucache) Set(key string, value interface{}) {

	if v, ok := lc.dict[key]; ok {
		lc.l.MoveToFront(v)
		v.Value.(*elem).v = value
		v.Value.(*elem).expire = 0
		return
	}

	if lc.l.Len() >= lc.cap {
		lc.DeleteOldest()
	}

	e := &elem{
		k:      key,
		v:      value,
		expire: 0,
	}
	node := lc.l.PushFront(e)
	lc.dict[key] = node

	return
}

// Get Get the value of a cached element by key. If key do not exist, this function will return nil and a error msg
// 		key:	The identity of an element
//		return:
//			value: 	the cached value, nil if key do not exist
// 			err:	error info, nil if value is not nil
func (lc *Lrucache) Get(key string) (value interface{}, err error) {
	if v, ok := lc.dict[key]; ok {
		lc.l.MoveToFront(v)
		return v.Value.(*elem).v, nil
	}
	return nil, errors.New("not found")
}

// Delete delete an element
func (lc *Lrucache) Delete(key string) error {
	if v, ok := lc.dict[key]; ok {
		lc.l.Remove(v)
		delete(lc.dict, key)
		return nil
	}
	return errors.New("not found")
}

// DeleteOldest delete the oldest element
func (lc *Lrucache) DeleteOldest() {
	oldest := lc.l.Back()
	if oldest != nil {
		oel := oldest.Value.(*elem)
		delete(lc.dict, oel.k)
		lc.l.Remove(oldest)
	}
}

// Keys get all cached unexpired keys.
func (lc *Lrucache) Keys() []string {
	var keys []string
	for k := range lc.dict {
		keys = append(keys, k)
	}
	return keys
}

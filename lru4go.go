package lru4go

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type elem struct {
	key interface{}
	data interface{}
	expireTime int64
	next *elem
	pre *elem
}

type lrucache struct {
	maxSize int
	elemCount int
	elemList map[interface{}]*elem
	first *elem
	last *elem
	mu sync.Mutex
}

// New 新建一个lrucache
func New(size int)(*lrucache, error) {
	newCache := new(lrucache)
	newCache.maxSize = size
	newCache.elemCount = 0
	return newCache, nil
}

// Set new or update an element
func (c *lrucache)Set(key interface{}, value interface{}, ttl...int) error {

	// 确保参数个数正确
	if len(ttl) > 1 {
		return errors.New("wrong para number, 2 or 3 expected but more than 3 received")
	}
	var ttlnum int64
	if len(ttl) == 1 {
		ttlnum = int64(ttl[0])
	} else {
		ttlnum = -1
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if e,ok := c.elemList[key]; ok {
		e.data = value
		if ttlnum == -1 {
			e.expireTime = ttlnum
		} else {
			e.expireTime = time.Now().Unix() + ttlnum
		}
		c.mvKeyToFirst(key)
	} else {
		if c.elemCount + 1 > c.maxSize {
			c.eliminationOldest()
		}
		newElem := &elem{
			key: key,
			data: value,
			expireTime: -1,
			pre: nil,
			next: c.first,
		}
		if ttlnum != -1 {
			newElem.expireTime = time.Now().Unix() + ttlnum
		}
		c.first = newElem
		c.elemList[key] = newElem

		c.elemCount++
	}
	return nil
}

// Get get an element by key
func (c *lrucache)Get(key interface{}) (value interface{}, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.elemList[key]; ok {
		if time.Now().Unix() > v.expireTime {
			// 如果过期了
			_ = c.Delete(key)
			return nil, errors.New("the key was expired")
		}
		c.mvKeyToFirst(key)
		return v.data, nil
	}
	return nil, errors.New("no value found")
}

// Delete delete an element
func (c *lrucache)Delete(key interface{}) error{
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.elemList[key]; ok {
		if v.pre == nil {
			// 当key是第一个元素时，清空元素列表，充值指针和元素计数
			c.elemList = make(map[interface{}]*elem)
			c.elemCount = 0
			c.last =  nil
			c.first = nil
			return nil
		} else if v.next == nil {
			// 当key不是第一个元素，但是是最后一个元素时,修改前一个元素的next指针并修改c.last指针
			v.pre.next = v.next
			c.last = v.pre
		} else {
			// 中间元素，修改前后指针
			v.pre.next = v.next
			v.next.pre = v.pre
		}
		delete(c.elemList, key)
		c.elemCount--
	}
	return errors.New(fmt.Sprintf("key %T do not exist", key))
}

// updateKeyPtr 更新对应key的指针，放到链表的第一个
func (c *lrucache)mvKeyToFirst(key interface{}) {
	elem := c.elemList[key]
	if elem.pre == nil {
		// 当key是第一个元素时，不做动作
		return
	} else if elem.next == nil {
		// 当key不是第一个元素，但是是最后一个元素时，提到第一个元素去
		elem.pre.next = nil

		c.last = elem.pre

		elem.pre = nil
		elem.next = c.first
		c.first = elem

	} else {
		elem.pre.next = elem.next
		elem.next.pre = elem.pre

		elem.next = c.first
		elem.pre = nil
		c.first = elem
	}
}

func (c *lrucache)  eliminationOldest() {
	if c.last == nil {
		return
	}
	if c.last.pre != nil {
		c.last.pre.next = nil
	}
	key := c.last.key
	c.last = c.last.pre
	delete(c.elemList, key)
}
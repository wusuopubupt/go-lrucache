package go_lrucache

import "errors"

type Key interface{}

type Value interface{}

type Node struct {
	next, prev *Node
	key        Key
	value      Value
}

// The cache interface
type Cache interface {
	Set(key Key, value Value)
	Get(key Key) (Value, bool)
	Del(key Key)
	Size() int
}

type LRUCache struct {
	Capacity   int
	kvMap      map[Key]*Node
	head, tail *Node
}

func New(capacity int) (*LRUCache, error) {
	if capacity == 0 {
		return nil, errors.New("Capacity can not be 0")
	}
	return &LRUCache{
		Capacity: capacity,
		kvMap:    make(map[Key]*Node),
	}, nil
}

/**
 * If map contains key, then update map and replace key to list head
 * Else set key to list head, check is map.size > capacity, if true, remove list tail
 */
func (lru *LRUCache) Set(key Key, value Value) {
	if lru.kvMap == nil {
		lru.kvMap = make(map[Key]*Node)
	}
	// key exists
	if node, ok := lru.kvMap[key]; ok {
		lru.remove(node)
		node.value = value
		lru.setHead(node)
		return
	}

	var node = &Node{nil, nil, key, value}
	lru.setHead(node)
	lru.kvMap[key] = node

	// remove oldest node
	if len(lru.kvMap) > lru.Capacity {
		lru.Del(lru.tail.key)
	}
}

func (lru *LRUCache) Get(key Key) (value Value, ok bool) {
	if node, ok := lru.kvMap[key]; ok {
		lru.remove(node)
		lru.setHead(node)
		return node.value, ok
	}
	return -1, false
}

func (lru *LRUCache) Del(key Key) {
	if node, hit := lru.kvMap[key]; hit {
		lru.RemoveNode(node)
	}
}

func (lru *LRUCache) Size() int {
	return lru.Capacity
}

func (lru *LRUCache) RemoveNode(node *Node) {
	// delete key from cache map
	delete(lru.kvMap, node.key)
	// remove node from cache list
	lru.remove(node)
}

func (lru *LRUCache) setHead(node *Node) {
	node.next = lru.head
	node.prev = nil
	if lru.head != nil {
		lru.head.prev = node
	}
	if lru.tail == nil {
		lru.tail = node
	}
	lru.head = node
}

func (lru *LRUCache) isHead(node *Node) bool {
	return node.prev == nil
}

func (lru *LRUCache) isTail(node *Node) bool {
	return node.next == nil
}

func (lru *LRUCache) remove(node *Node) {
	if lru.isHead(node) {
		lru.head = node.next
	} else {
		node.prev.next = node.next
	}

	if lru.isTail(node) {
		lru.tail = node.prev
	} else {
		node.next.prev = node.prev
	}

}

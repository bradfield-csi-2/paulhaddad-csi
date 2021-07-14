package main

import (
	"math/rand"
	"time"
)

const maxLevel = 10
const p = 0.5

type skipNode struct {
	item    Item
	forward []*skipNode
}

type skipListOC struct {
	header *skipNode
	tail   *skipNode
	level  int
}

func newSkipListOC() *skipListOC {
	header := &skipNode{
		forward: make([]*skipNode, maxLevel),
	}
	tail := &skipNode{
		Item{"zzzzzzz", "zzzzzzzz"},
		nil,
	}
	header.forward[0] = tail
	return &skipListOC{header, tail, 1}
}

func (o *skipListOC) Get(key string) (string, bool) {
	x := o.header

	for i := o.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].item.Key < key {
			x = x.forward[i]
		}
	}

	x = x.forward[0]
	if x.item.Key != key {
		return "", false
	}

	return x.item.Value, true
}

func (o *skipListOC) Put(key, value string) bool {
	update := make([]*skipNode, maxLevel)
	x := o.header

	for i := o.level - 1; i >= 0; i-- {

		for x.forward[i] != nil && x.forward[i].item.Key < key {
			x = x.forward[i]
		}
		update[i] = x
	}

	x = x.forward[0]

	if x.item.Key == key {
		x.item.Value = value
		return false
	}

	level := randLevel()
	if level > o.level {
		for i := o.level - 1; i <= level-1; i++ {
			update[i] = o.header
		}

		o.level = level
	}

	x = &skipNode{
		Item{key, value},
		make([]*skipNode, maxLevel),
	}

	for i := 0; i <= level-1; i++ {
		x.forward[i] = update[i].forward[i]
		update[i].forward[i] = x
	}

	return true
}

func (o *skipListOC) Delete(key string) bool {
	update := make([]*skipNode, maxLevel)
	x := o.header

	for i := o.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].item.Key < key {
			x = x.forward[i]
		}
		update[i] = x
	}

	x = x.forward[0]

	if x.item.Key == key {
		update[0].forward[0] = x.forward[0]
		return true
	}

	return false
}

func randLevel() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	level := 0
	for r.Float32() < p && level < maxLevel-1 {
		level++
	}

	return level
}

func (o *skipListOC) RangeScan(startKey, endKey string) Iterator {
	return &skipListOCIterator{}
}

type skipListOCIterator struct {
}

func (iter *skipListOCIterator) Next() {
}

func (iter *skipListOCIterator) Valid() bool {
	return false
}

func (iter *skipListOCIterator) Key() string {
	return ""
}

func (iter *skipListOCIterator) Value() string {
	return ""
}

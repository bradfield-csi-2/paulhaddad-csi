package main

type node struct {
	item Item
	next *node
}

type skipListOC struct {
	head *node
}

func newSkipListOC() *skipListOC {
	return &skipListOC{}
}

// The second return value will be `false` when the `key` hasn't been
// associated with any value.
func (o *skipListOC) Get(key string) (string, bool) {
	curNode := o.head
	if curNode == nil {
		return "", false
	}

	for curNode != nil && curNode.item.Key < key {
		curNode = curNode.next
	}

	if curNode.item.Key == key {
		return curNode.item.Value, true
	}

	return "", false
}

// Put should return `true` if a new key was added, and `false` if an
// existing key had its value updated.
func (o *skipListOC) Put(key, value string) bool {
	curNode := o.head

	if curNode == nil {
		item := Item{key, value}
		o.head = &node{item: item, next: nil}
		return true
	}

	prevNode := curNode
	curNode = curNode.next

	for curNode != nil && curNode.item.Key < key {
		prevNode, curNode = curNode, curNode.next
	}

	// add new node to end of list
	if curNode == nil {
		item := Item{key, value}
		prevNode.next = &node{item: item, next: curNode}
		return true
	}

	// updating existing key
	if curNode.item.Key == key {
		curNode.item.Value = value
		return false
	}

	// create new key
	item := Item{key, value}
	prevNode.next = &node{item: item, next: curNode}
	return true
}

// Delete should return whether or not the key was actually deleted, i.e.
// it should return `true` if the key existed before deletion.
func (o *skipListOC) Delete(key string) bool {
	prevNode := o.head
	if prevNode == nil {
		return false
	}

	curNode := prevNode.next

	if prevNode.item.Key == key {
		o.head = curNode
		return true
	}

	for curNode != nil && curNode.item.Key < key {
		prevNode, curNode = curNode, curNode.next
	}

	if curNode.item.Key == key {
		prevNode.next = curNode.next
		return true
	}

	return false
}

// startKey and endKey are inclusive.
func (o *skipListOC) RangeScan(startKey, endKey string) Iterator {
	node := o.head

	for node != nil && node.item.Key < startKey {
		node = node.next
	}

	return &skipListOCIterator{o, node, startKey, endKey}
}

type skipListOCIterator struct {
	o                *skipListOC
	node             *node
	startKey, endKey string
}

func (iter *skipListOCIterator) Next() {
	iter.node = iter.node.next
}

func (iter *skipListOCIterator) Valid() bool {
	return iter.node != nil && iter.node.item.Key <= iter.endKey
}

func (iter *skipListOCIterator) Key() string {
	return iter.node.item.Key
}

func (iter *skipListOCIterator) Value() string {
	return iter.node.item.Value
}

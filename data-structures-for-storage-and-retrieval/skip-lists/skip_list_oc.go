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

	for curNode != nil {
		if curNode.item.Key == key {
			return curNode.item.Value, true
		}

		curNode = curNode.next
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

	for curNode.next != nil {
		if curNode.item.Key == key {
			curNode.item.Value = value
			return false
		}
		curNode = curNode.next
	}

	item := Item{key, value}
	curNode.next = &node{item: item, next: nil}
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

	for curNode != nil {
		if curNode.item.Key == key {
			prevNode.next = curNode.next
			return true
		}
		prevNode = curNode
		curNode = curNode.next
	}

	return false
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

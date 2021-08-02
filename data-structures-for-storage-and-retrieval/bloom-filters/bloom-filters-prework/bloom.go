package main

import (
	"encoding/binary"
	"hash"
)

const bits = 8

type bloomFilter interface {
	add(item string)

	// `false` means the item is definitely not in the set
	// `true` means the item might be in the set
	maybeContains(item string) bool

	// Number of bytes used in any underlying storage
	memoryUsage() int
}

type realBloomFilter struct {
	data      []uint64
	numHashes int
	hashFn    func() hash.Hash32
	m         int
}

func newRealBloomFilter(size, k int, fn func() hash.Hash32) *realBloomFilter {
	elements := size / 8

	return &realBloomFilter{
		data:      make([]uint64, elements),
		numHashes: k,
		hashFn:    fn,
		m:         elements * 64,
	}
}

func (b *realBloomFilter) add(item string) {
	for i := 0; i < b.numHashes; i++ {
		block, offset := b.hashLocation(string(i) + item)
		b.data[block] |= uint64(offset)
	}
}

func (b *realBloomFilter) maybeContains(item string) bool {
	// Technically, any item "might" be in the set
	for i := 0; i < b.numHashes; i++ {
		block, offset := b.hashLocation(string(i) + item)
		if b.data[block]&uint64(offset) == 0 {
			return false
		}
	}

	return true
}

func (b *realBloomFilter) hashLocation(key string) (int, int) {
	hash := b.hashFn()
	hash.Write([]byte(key))
	x, _ := binary.Uvarint(hash.Sum(nil))
	x = x % uint64(b.m)

	block, index := x/64, x%64
	offset := 1 << index

	return int(block), int(offset)
}

func (b *realBloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}

type trivialBloomFilter struct {
	data []uint64
}

func newTrivialBloomFilter() *trivialBloomFilter {
	return &trivialBloomFilter{
		data: make([]uint64, 1000),
	}
}

func (b *trivialBloomFilter) add(item string) {
	// Do nothing
}

func (b *trivialBloomFilter) maybeContains(item string) bool {
	// Technically, any item "might" be in the set
	return true
}

func (b *trivialBloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}

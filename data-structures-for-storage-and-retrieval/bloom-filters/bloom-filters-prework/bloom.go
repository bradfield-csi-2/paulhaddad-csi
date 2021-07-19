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
	m         uint64
	data      []uint64
	numHashes int
	hashFn    func() hash.Hash64
}

func newRealBloomFilter(m uint64, k int, fn func() hash.Hash64) *realBloomFilter {
	return &realBloomFilter{
		m:         m,
		data:      make([]uint64, m/64),
		numHashes: k,
		hashFn:    fn,
	}
}

func (b *realBloomFilter) add(item string) {
	for i := 0; i < b.numHashes; i++ {
		block, offset := b.hashLocation(string(i) + item)
		b.data[block] = offset | b.data[block]
	}
}

func (b *realBloomFilter) maybeContains(item string) bool {
	// Technically, any item "might" be in the set
	for i := 0; i < b.numHashes; i++ {
		block, offset := b.hashLocation(string(i) + item)
		if offset&b.data[block] == 0 {
			return false
		}
	}

	return true
}

func (b *realBloomFilter) hashLocation(key string) (uint64, uint64) {
	hash := b.hashFn()
	hash.Write([]byte(key))
	x, _ := binary.Uvarint(hash.Sum(nil))
	x = x % (b.m / bits)

	block, index := x/bits, x%8
	offset := uint64(1 << index)

	return block, offset
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

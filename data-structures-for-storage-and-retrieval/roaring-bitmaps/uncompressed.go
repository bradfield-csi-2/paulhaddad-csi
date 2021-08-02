package bitmap

import (
	"math"
)

const wordSize = 64

type uncompressedBitmap struct {
	data []uint64
}

func newUncompressedBitmap() *uncompressedBitmap {
	return &uncompressedBitmap{
		data: make([]uint64, math.MaxUint32/wordSize),
	}
}

func (b *uncompressedBitmap) Get(x uint32) bool {
	index := x / 64
	offset := x % 64
	return b.data[index]&(1<<offset) > 0
}

func (b *uncompressedBitmap) Set(x uint32) {
	index := x / 64
	offset := x % 64
	b.data[index] |= 1 << offset
}

func (b *uncompressedBitmap) Union(other *uncompressedBitmap) *uncompressedBitmap {
	data := make([]uint64, math.MaxUint32/wordSize)

	for i := range data {
		data[i] = b.data[i] | other.data[i]
	}

	return &uncompressedBitmap{
		data: data,
	}
}

func (b *uncompressedBitmap) Intersect(other *uncompressedBitmap) *uncompressedBitmap {
	data := make([]uint64, math.MaxUint32/wordSize)

	for i := range data {
		data[i] = b.data[i] & other.data[i]
	}
	return &uncompressedBitmap{
		data: data,
	}
}

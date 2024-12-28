package main

import (
	"fmt"
)

type BitMap struct {
	data []uint64
}

func NewBitMap(size int) *BitMap {
	return &BitMap{
		data: make([]uint64, (size+63)/64),
	}
}

func (b *BitMap) Set(index int) {
	if index < 0 {
		return
	}
	i := index / 64
	offset := index % 64
	b.data[i] |= (1 << offset)
}

func (b *BitMap) Test(index int) bool {
	if index < 0 {
		return false
	}
	i := index / 64
	offset := index % 64
	return (b.data[i] & (1 << offset)) != 0
}

func main() {
	bm := NewBitMap(100)
	bm.Set(5)
	bm.Set(10)
	fmt.Println(bm.Test(5))
	fmt.Println(bm.Test(8))
}

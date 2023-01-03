package main

import (
	"math"
)

var magicTable = []uint64{
	0, 1, 2, 53, 3, 7, 54, 27, 4, 38, 41, 8, 34, 55, 48, 28,
	62, 5, 39, 46, 44, 42, 22, 9, 24, 35, 59, 56, 49, 18, 29, 11,
	63, 52, 6, 26, 37, 40, 33, 47, 61, 45, 43, 21, 23, 58, 17, 10,
	51, 25, 36, 32, 60, 20, 57, 16, 50, 31, 19, 15, 30, 14, 13, 12,
}

// Log2 calculates log2 of a uint64 value.
// Ported from
// https://github.com/h2o/h2o/blob/0f08b675c8244fc4552a93e9b35271ecf5e0f8fa/deps/libgkc/gkc.c#L109-L127
func Log2(x uint64) uint64 {
	const debruijnMagic = uint64(0x022fdd63cc95386d)

	x |= (x >> 1)
	x |= (x >> 2)
	x |= (x >> 4)
	x |= (x >> 8)
	x |= (x >> 16)
	x |= (x >> 32)
	return (magicTable[((x & ^(x>>1))*debruijnMagic)>>58])
}

var table = []uint64{
	0, 58, 1, 59, 47, 53, 2, 60, 39, 48, 27, 54, 33, 42, 3, 61,
	51, 37, 40, 49, 18, 28, 20, 55, 30, 34, 11, 43, 14, 22, 4, 62,
	57, 46, 52, 38, 26, 32, 41, 50, 36, 17, 19, 29, 10, 13, 21, 56,
	45, 25, 31, 35, 16, 9, 12, 44, 24, 15, 8, 23, 7, 6, 5, 63,
}

// Ported from https://stackoverflow.com/a/23000588/1391518
func Log2ByAvernar(n uint64) uint64 {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return table[(n*0x03f6eaf2cd271461)>>58]
}

func buildTable(c uint64) []uint8 {
	v := make([]uint8, 64)
	for i := uint64(0); i < 64; i++ {
		x := uint64(1) << i
		x |= x - 1
		v[(x*c)>>58] = uint8(i)
	}
	return v
}

var u8Table = []uint8{
	0, 58, 1, 59, 47, 53, 2, 60, 39, 48, 27, 54, 33, 42, 3, 61,
	51, 37, 40, 49, 18, 28, 20, 55, 30, 34, 11, 43, 14, 22, 4, 62,
	57, 46, 52, 38, 26, 32, 41, 50, 36, 17, 19, 29, 10, 13, 21, 56,
	45, 25, 31, 35, 16, 9, 12, 44, 24, 15, 8, 23, 7, 6, 5, 63,
}

// Ported from https://stackoverflow.com/a/23000588/1391518
func Log2ByAvernarU8(n uint64) uint64 {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return uint64(u8Table[(n*0x03f6eaf2cd271461)>>58])
}

func Log2ByStdlib(n uint64) uint64 {
	return uint64(math.Floor(math.Log2(float64(n))))
}

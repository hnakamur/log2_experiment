package main

import (
	"math"
	"math/bits"
)

// ILog2 calculates log2 of a uint64 value.
// Ported from
// https://github.com/h2o/h2o/pull/3177/files
func ILog2(x uint64) int64 {
	return 63 - int64(bits.LeadingZeros64(x))
}

// ILog2B calculates log2 of a uint64 value.
// Ported from
// https://twitter.com/herumi/status/1610248792254844929
func ILog2B(x uint64) int64 {
	f := float64(x)
	v := math.Float64bits(f)
	return int64(v>>52) - 1023
}

var table = []int64{
	0, 58, 1, 59, 47, 53, 2, 60, 39, 48, 27, 54, 33, 42, 3, 61,
	51, 37, 40, 49, 18, 28, 20, 55, 30, 34, 11, 43, 14, 22, 4, 62,
	57, 46, 52, 38, 26, 32, 41, 50, 36, 17, 19, 29, 10, 13, 21, 56,
	45, 25, 31, 35, 16, 9, 12, 44, 24, 15, 8, 23, 7, 6, 5, 63,
}

// Ported from https://stackoverflow.com/a/23000588/1391518
func Log2ByAvernar(n uint64) int64 {
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
func Log2ByAvernarU8(n uint64) int64 {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	return int64(u8Table[(n*0x03f6eaf2cd271461)>>58])
}

func Log2Double(n uint64) int64 {
	return int64(math.Log2(float64(n)))
}

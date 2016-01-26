package mempool

// http://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2
// return the power of 2 which is larger or equal the target
// 000 -> 1 (001)
// 001 -> 1 (001)
// 010 -> 2 (010)
// 011 -> 3 (100)
// 100 -> 3 (100)
// 101 -> 4 (1000)
// 110 -> 4 (1000)
// 111 -> 4 (1000)
func LargerPowerOf2(v uint64) (r byte) {
	if v <= 1 {
		return 0
	}

	v--
	r = largestSet(v)
	r++

	return r
}

func largestSet(v64 uint64) (r byte) {

	if v64&0xffffffff00000000 != 0 {
		return largestSet32(uint32(v64>>32)) + 32
	}
	return largestSet32(uint32(v64))

}

var multiplyDeBruijnBitPosition = [32]byte{
	0, 9, 1, 10, 13, 21, 2, 29, 11, 14, 16, 18, 22, 25, 3, 30,
	8, 12, 20, 28, 15, 17, 24, 7, 19, 27, 23, 6, 26, 5, 4, 31,
}

func largestSet32(v uint32) (r byte) {
	// http://graphics.stanford.edu/~seander/bithacks.html#IntegerLogDeBruijn
	// Returns the index of the largest set bit.
	v |= v >> 1 // first round down to one less than a power of 2
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16

	r = multiplyDeBruijnBitPosition[(v*0x07c4acdd)>>27]
	return
}

package skein

type block256_slow struct {
	state [4]uint64
	key   [5]uint64 // The user-defined key is the first 4 values, we add the 5th.
	tweak [3]uint64 // The user-defined tweak is the first 2 values, we add the 3rd.
}

func (b *block256_slow) encrypt() {
	b.key[4] = c240 ^ b.key[0] ^ b.key[1] ^ b.key[2] ^ b.key[3]
	b.tweak[2] = b.tweak[0] ^ b.tweak[1]

	var r int
	// Our loop will run four rounds at a time, since the key-schedule only comes into play on every
	// fourth round.
	for r = 0; r < 18; r++ {
		b.state[0] += b.key[r%5]
		b.state[1] += b.key[(r+1)%5] + b.tweak[r%3]
		b.state[2] += b.key[(r+2)%5] + b.tweak[(r+1)%3]
		b.state[3] += b.key[(r+3)%5] + uint64(r)
		b.mix(r * 4)
		b.permute()
		b.mix(r*4 + 1)
		b.permute()
		b.mix(r*4 + 2)
		b.permute()
		b.mix(r*4 + 3)
		b.permute()
	}
	b.state[0] += b.key[r%5]
	b.state[1] += b.key[(r+1)%5] + b.tweak[r%3]
	b.state[2] += b.key[(r+2)%5] + b.tweak[(r+1)%3]
	b.state[3] += b.key[(r+3)%5] + uint64(r)
}

func (b *block256_slow) decrypt() {
	b.key[4] = c240 ^ b.key[0] ^ b.key[1] ^ b.key[2] ^ b.key[3]
	b.tweak[2] = b.tweak[0] ^ b.tweak[1]

	r := 18
	b.state[0] -= b.key[r%5]
	b.state[1] -= b.key[(r+1)%5] + b.tweak[r%3]
	b.state[2] -= b.key[(r+2)%5] + b.tweak[(r+1)%3]
	b.state[3] -= b.key[(r+3)%5] + uint64(r)
	// Our loop will run four rounds at a time, since the key-schedule only comes into play on every
	// fourth round.
	for r := 17; r >= 0; r-- {
		b.unpermute()
		b.unmix(r*4 + 3)
		b.unpermute()
		b.unmix(r*4 + 2)
		b.unpermute()
		b.unmix(r*4 + 1)
		b.unpermute()
		b.unmix(r * 4)

		b.state[0] -= b.key[r%5]
		b.state[1] -= b.key[(r+1)%5] + b.tweak[r%3]
		b.state[2] -= b.key[(r+2)%5] + b.tweak[(r+1)%3]
		b.state[3] -= b.key[(r+3)%5] + uint64(r)
	}
}

func (b *block256_slow) mix(d int) {
	b.state[0], b.state[1] = mix256_slow(d, 0, b.state[0], b.state[1])
	b.state[2], b.state[3] = mix256_slow(d, 1, b.state[2], b.state[3])
}

func (b *block256_slow) unmix(d int) {
	b.state[0], b.state[1] = unmix256_slow(d, 0, b.state[0], b.state[1])
	b.state[2], b.state[3] = unmix256_slow(d, 1, b.state[2], b.state[3])
}

func (b *block256_slow) permute() {
	// The permutation for 256 just swaps positions 1 and 3.
	// i.e. perm = [4]int{0, 3, 2, 1}
	b.state[1], b.state[3] = b.state[3], b.state[1]
}

func (b *block256_slow) unpermute() {
	// The permutation for 256 just swaps positions 1 and 3.
	// i.e. perm = [4]int{0, 3, 2, 1}
	b.state[1], b.state[3] = b.state[3], b.state[1]
}

func mix256_slow(d, j int, x0, x1 uint64) (y0, y1 uint64) {
	r := tf256Rots_slow[d&0x07][j]
	y0 = x0 + x1
	y1 = (x1 << r) | (x1 >> (64 - r))
	y1 = y1 ^ y0
	return
}

func unmix256_slow(d, j int, y0, y1 uint64) (x0, x1 uint64) {
	r := tf256Rots_slow[d&0x07][j]
	x1 = y1 ^ y0
	x1 = (x1 >> r) | (x1 << (64 - r))
	x0 = y0 - x1
	return
}

var tf256Rots_slow = [8][2]uint{
	{14, 16},
	{52, 57},
	{23, 40},
	{5, 37},
	{25, 33},
	{46, 12},
	{58, 22},
	{32, 32},
}

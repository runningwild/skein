// +build amd64

package skein

type block256 struct {
	state [4]uint64
	key   [5]uint64 // The user-defined key is the first 4 values, we add the 5th.
	tweak [3]uint64 // The user-defined tweak is the first 2 values, we add the 3rd.
}

func (b *block256) encrypt() {
	encryptBlock256(&b.state, &b.key, &b.tweak)
}
func encryptBlock256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64) {
	key[4] = c240 ^ key[0] ^ key[1] ^ key[2] ^ key[3]
	tweak[2] = tweak[0] ^ tweak[1]

	var d int
	for d = 0; d < 18; d++ {
		keySched(d, state, key, tweak, &keyShifts, &tweakShifts)
		fourRoundsA(state)
		d++
		keySched(d, state, key, tweak, &keyShifts, &tweakShifts)
		fourRoundsB(state)
	}
	state[0] += key[d%5]
	state[1] += key[(d+1)%5] + tweak[d%3]
	state[2] += key[(d+2)%5] + tweak[(d+1)%3]
	state[3] += key[(d+3)%5] + uint64(d)
}

var keyShifts = [85]int{
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
	0, 1, 2, 3, 4,
}

var tweakShifts = [85]int{
	0, 1, 2, 0, 1, 2, 0, 1, 2,
	0, 1, 2, 0, 1, 2, 0, 1, 2,
	0, 1, 2, 0, 1, 2, 0, 1, 2,
	0, 1, 2, 0, 1, 2, 0, 1, 2,
	0, 1, 2, 0, 1, 2, 0, 1, 2,
	0, 1, 2, 0, 1, 2, 0, 1, 2,
	0, 1, 2, 0, 1, 2, 0, 1, 2,
	0, 1, 2, 0, 1, 2, 0, 1, 2,
	0, 1, 2, 0, 1, 2, 0, 1, 2,
	0, 1, 2, 0,
}

func keySched(d int, state *[4]uint64, key *[5]uint64, tweak *[3]uint64, keyShifts, tweakShifts *[85]int) {
	state[0] += key[keyShifts[d]]
	state[1] += key[keyShifts[d+1]] + tweak[tweakShifts[d]]
	state[2] += key[keyShifts[d+2]] + tweak[tweakShifts[d+1]]
	state[3] += key[keyShifts[d+3]] + uint64(d)
}
func eightRounds(d int, state *[4]uint64, key *[5]uint64, tweak *[3]uint64, keyShifts, tweakShifts *[85]int) //{}
func fourRoundsA(state *[4]uint64)
func fourRoundsB(state *[4]uint64)
func fourRoundsA_slow(state *[4]uint64) {
	state[0] += state[1]
	state[1] = ((state[1] << 14) | (state[1] >> (64 - 14))) ^ state[0]
	state[2] += state[3]
	state[3] = ((state[3] << 16) | (state[3] >> (64 - 16))) ^ state[2]
	state[0] += state[3]
	state[3] = ((state[3] << 52) | (state[3] >> (64 - 52))) ^ state[0]
	state[2] += state[1]
	state[1] = ((state[1] << 57) | (state[1] >> (64 - 57))) ^ state[2]

	state[0] += state[1]
	state[1] = ((state[1] << 23) | (state[1] >> (64 - 23))) ^ state[0]
	state[2] += state[3]
	state[3] = ((state[3] << 40) | (state[3] >> (64 - 40))) ^ state[2]

	state[0] += state[3]
	state[3] = ((state[3] << 5) | (state[3] >> (64 - 5))) ^ state[0]
	state[2] += state[1]
	state[1] = ((state[1] << 37) | (state[1] >> (64 - 37))) ^ state[2]
}

func fourRoundsB_slow(state *[4]uint64) {
	state[0] += state[1]
	state[1] = ((state[1] << 25) | (state[1] >> (64 - 25))) ^ state[0]
	state[2] += state[3]
	state[3] = ((state[3] << 33) | (state[3] >> (64 - 33))) ^ state[2]

	state[0] += state[3]
	state[3] = ((state[3] << 46) | (state[3] >> (64 - 46))) ^ state[0]
	state[2] += state[1]
	state[1] = ((state[1] << 12) | (state[1] >> (64 - 12))) ^ state[2]

	state[0] += state[1]
	state[1] = ((state[1] << 58) | (state[1] >> (64 - 58))) ^ state[0]
	state[2] += state[3]
	state[3] = ((state[3] << 22) | (state[3] >> (64 - 22))) ^ state[2]

	state[0] += state[3]
	state[3] = ((state[3] << 32) | (state[3] >> (64 - 32))) ^ state[0]
	state[2] += state[1]
	state[1] = ((state[1] << 32) | (state[1] >> (64 - 32))) ^ state[2]
}

func (b *block256) decrypt() {
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

// func (b *block256) mix(d int) {
// 	b.state[0], b.state[1] = mix256(d, 0, b.state[0], b.state[1])
// 	b.state[2], b.state[3] = mix256(d, 1, b.state[2], b.state[3])
// }

func (b *block256) unmix(d int) {
	b.state[0], b.state[1] = unmix256(d, 0, b.state[0], b.state[1])
	b.state[2], b.state[3] = unmix256(d, 1, b.state[2], b.state[3])
}

func (b *block256) permute() {
	// The permutation for 256 just swaps positions 1 and 3.
	// i.e. perm = [4]int{0, 3, 2, 1}
	b.state[1], b.state[3] = b.state[3], b.state[1]
}

func (b *block256) unpermute() {
	// The permutation for 256 just swaps positions 1 and 3.
	// i.e. perm = [4]int{0, 3, 2, 1}
	b.state[1], b.state[3] = b.state[3], b.state[1]
}

func mix256(rot uint, x0, x1 uint64) (y0, y1 uint64) { return }
func mix256_bad(rot uint, x0, x1 uint64) (y0, y1 uint64) {
	y0 = x0 + x1
	y1 = (x1 << rot) | (x1 >> (64 - rot))
	y1 = y1 ^ y0
	return
}

func unmix256(d, j int, y0, y1 uint64) (x0, x1 uint64) {
	r := tf256Rots[d&0x07][j]
	x1 = y1 ^ y0
	x1 = (x1 >> r) | (x1 << (64 - r))
	x0 = y0 - x1
	return
}

var tf256Rots = [8][2]uint{
	{14, 16},
	{52, 57},
	{23, 40},
	{5, 37},
	{25, 33},
	{46, 12},
	{58, 22},
	{32, 32},
}

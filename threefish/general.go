package threefish

import (
	"fmt"

	"github.com/runningwild/skein/convert"
)

type General struct {
	mix          [][]uint
	perm         []int
	inversePerm  []int
	numRounds    int
	keySchedule  int
	keyExtension uint64
}

func NewGeneral(mix [][]uint, perm []int, numRounds int, keySchedule int, keyExtension uint64) (*General, error) {
	period := findPeriod(perm)
	if len(mix)%period != 0 {
		return nil, fmt.Errorf("the period of the permutation (%d) must divide the number of mixes (%d)", period, len(mix))
	}
	for _, v := range mix {
		if len(v) != len(perm)/2 {
			return nil, fmt.Errorf("every element of mix must be len(perm)/2")
		}
	}
	if numRounds%keySchedule != 0 {
		return nil, fmt.Errorf("the key schedule must divide the number of rounds")
	}
	if numRounds%len(mix) != 0 {
		return nil, fmt.Errorf("the number of mixes must divide the number of rounds")
	}

	inverse := make([]int, len(perm))
	for i := range inverse {
		inverse[perm[i]] = i
	}

	// Easier to use the inverse when permuting them.
	perm, inverse = inverse, perm

	return &General{
		mix:          mix,
		perm:         perm,
		inversePerm:  inverse,
		numRounds:    numRounds,
		keySchedule:  keySchedule,
		keyExtension: keyExtension,
	}, nil
}

func (g *General) MakeCipher(key []byte) (*GeneralCipher, error) {
	if len(key) != len(g.perm)*8 {
		return nil, fmt.Errorf("key must be %d bytes", len(g.perm)*8)
	}
	c := &GeneralCipher{
		gen:       g,
		key:       make([]uint64, len(key)/8+1),
		blockSize: len(key),
	}
	copy(c.key, convert.InplaceBytesToUInt64(key))
	for i := range c.key {
		c.key[len(c.key)-1] ^= c.key[i]
	}
	c.key[len(c.key)-1] ^= c.gen.keyExtension
	return c, nil
}

type GeneralCipher struct {
	gen       *General
	key       []uint64
	blockSize int
}

func (c *GeneralCipher) Encrypt(dst, src []byte) {
	if len(dst) != c.blockSize || len(src) != c.blockSize {
		panic(fmt.Sprintf("dst and src were of size %d and %d, but blocksize is %d", len(dst), len(src), c.blockSize))
	}
	copy(dst, src)
	block := convert.InplaceBytesToUInt64(dst)
	sched := 0
	for r := 0; r < c.gen.numRounds; r++ {
		if r%c.gen.keySchedule == 0 {
			c.keySchedule(sched, block)
			sched++
		}
		c.mix(r, block)
		c.permute(block)
	}
	c.keySchedule(sched, block)
}

func (c *GeneralCipher) Decrypt(dst, src []byte) {
	if len(dst) != c.blockSize || len(src) != c.blockSize {
		panic(fmt.Sprintf("dst and src were of size %d and %d, but blocksize is %d", len(dst), len(src), c.blockSize))
	}
	copy(dst, src)
	block := convert.InplaceBytesToUInt64(dst)
	sched := c.gen.numRounds / c.gen.keySchedule
	c.keyUnschedule(sched, block)
	sched--
	for r := c.gen.numRounds - 1; r >= 0; r-- {
		c.unpermute(block)
		c.unmix(r, block)
		if r%c.gen.keySchedule == 0 {
			c.keyUnschedule(sched, block)
			sched--
		}
	}
}

func (c *GeneralCipher) keySchedule(round int, block []uint64) {
	for i := range block {
		block[i] += c.key[(i+round)%len(c.key)]
	}
	block[len(block)-2] += 0             // tweak
	block[len(block)-1] += 0             // tweak
	block[len(block)-1] += uint64(round) // tweak
}
func (c *GeneralCipher) mix(round int, block []uint64) {
	for i, rot := range c.gen.mix[round%len(c.gen.mix)] {
		a := i * 2
		b := a + 1
		block[a] += block[b]
		block[b] = ((block[b] << rot) | (block[b] >> (64 - rot))) ^ block[a]
	}
}
func (c *GeneralCipher) permute(block []uint64) {
	used := make([]bool, len(c.gen.perm))
	for i := range c.gen.perm {
		if used[i] {
			continue
		}
		buf := block[i]
		for j := c.gen.perm[i]; j != i; j = c.gen.perm[j] {
			used[j] = true
			block[j], buf = buf, block[j]
		}
		block[i] = buf
	}
}

func (c *GeneralCipher) keyUnschedule(round int, block []uint64) {
	for i := range block {
		block[i] -= c.key[(i+round)%len(c.key)]
	}
	block[len(block)-2] -= 0             // tweak
	block[len(block)-1] -= 0             // tweak
	block[len(block)-1] -= uint64(round) // tweak
}
func (c *GeneralCipher) unmix(round int, block []uint64) {
	for i, rot := range c.gen.mix[round%len(c.gen.mix)] {
		a := i * 2
		b := a + 1
		block[b] ^= block[a]
		block[b] = (block[b] >> rot) | (block[b] << (64 - rot))
		block[a] -= block[b]
	}
}
func (c *GeneralCipher) unpermute(block []uint64) {
	used := make([]bool, len(c.gen.inversePerm))
	for i := range c.gen.inversePerm {
		if used[i] {
			continue
		}
		buf := block[i]
		for j := c.gen.inversePerm[i]; j != i; j = c.gen.inversePerm[j] {
			used[j] = true
			block[j], buf = buf, block[j]
		}
		block[i] = buf
	}
}

func findPeriod(perm []int) int {
	used := make([]bool, len(perm))
	var cycles []int
	for i := range perm {
		if used[i] {
			continue
		}
		used[i] = true
		length := 1
		for j := perm[i]; j != i; j = perm[j] {
			used[j] = true
		}
		cycles = append(cycles, length)
	}
	period := 1
	for _, c := range cycles {
		period *= c
	}
	return period
}

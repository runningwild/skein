package threefish

import (
	"unsafe"

	"github.com/runningwild/skein/convert"
	"github.com/runningwild/skein/types"
)

func init() {
	jfishRegister("avx2", makeJFishAVX2)
}

func makeJFishAVX2() types.JFish {
	// 32 uint64s for the states
	// 32 uint64s for the keys
	// 8 uint64s for the tweaks
	v := make32ByteAligned(32 + 32 + 8)
	states := v[0:32]
	keys := v[32:64]
	tweaks := v[64:70]

	j := jFishAVX2{
		statePtr: &states[0],
		keyPtr:   &keys[0],
		tweakPtr: &tweaks[0],
	}
	for i := range j.state {
		j.state[i] = convert.InplaceUint64ToBytes(states[i*8 : (i+1)*8])
	}
	for i := range j.key {
		j.key[i] = convert.InplaceUint64ToBytes(keys[i*2 : (i+1)*2])
	}
	for i := range j.tweak {
		j.tweak[i] = convert.InplaceUint64ToBytes(tweaks[i*2 : (i+1)*2])
	}

	return &j
}

type jFishAVX2 struct {
	state [4][]byte
	key   [4][]byte
	tweak [4][]byte

	statePtr *uint64
	keyPtr   *uint64
	tweakPtr *uint64
}

func (j *jFishAVX2) NumLanes() int {
	return 4
}

func (j *jFishAVX2) State(lane int) []byte {
	return j.state[lane]
}

func (j *jFishAVX2) Tweak(lane int) []byte {
	return j.tweak[lane]
}

func (j *jFishAVX2) Encrypt(states, keys, tweaks [][]byte) {
	if len(states) != len(keys) || len(keys) != len(tweaks) || len(states) > j.NumLanes() {
		panic("invalid lane arrangement")
	}
	for i := range states {
		copy(j.state[i], states[i])
		copy(j.key[i], keys[i])
		copy(j.tweak[i], tweaks[i])
	}
	threefish512_avx2(j.statePtr, j.keyPtr, j.tweakPtr)
	for i := range states {
		copy(states[i], j.state[i])
	}
}

func make32ByteAligned(n int) []uint64 {
	v := make([]uint64, n+4)
	for i := range v {
		if uintptr(unsafe.Pointer(&v[i]))%32 == 0 {
			return v[i : i+n]
		}
	}
	panic("X")
}

func threefish512_avx2(state *uint64, key *uint64, tweak *uint64)
func transpose4x8()

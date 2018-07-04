package threefish

import (
	"unsafe"

	"github.com/runningwild/skein/convert"
)

func init() {
	jfishRegister("avx2", makeJFishAVX2)
}

func makeJFishAVX2(key [64]byte) JFish {
	// 32 uint64s for the state
	// 8 uint64s for the shared key
	// 8 uint64s for the tweaks
	v := make32ByteAligned(32 + 8 + 8)
	states := v[0:32]
	tweaks := v[32:40]

	j := jFishAVX2{
		statePtr:  &states[0],
		tweakPtr:  &tweaks[0],
		sharedKey: v[40:48],
	}
	for i := range j.state {
		j.state[i] = convert.InplaceUint64ToBytes(states[i*8 : (i+1)*8])
	}
	for i := range j.tweak {
		j.tweak[i] = convert.InplaceUint64ToBytes(tweaks[i*2 : (i+1)*2])
	}
	copy(j.sharedKey[:], convert.InplaceBytesToUint64(key[:]))

	return &j
}

type jFishAVX2 struct {
	state [4][]byte
	tweak [4][]byte

	sharedKey []uint64

	statePtr *uint64
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

func (j *jFishAVX2) Encrypt() {
	threefish512_avx2(j.statePtr, &j.sharedKey[0], j.tweakPtr)
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

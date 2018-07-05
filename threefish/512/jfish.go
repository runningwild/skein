package threefish

import (
	"fmt"

	"github.com/intel-go/cpuid"
	"github.com/runningwild/skein/convert"
)

type JFish interface {
	// Returns the number of lanes this JFish object operates on.
	NumLanes() int

	// Returns the state for the specified lane.
	State(lane int) []byte

	// Returns the tweaks for the specified lane.
	Tweak(lane int) []byte

	// Encrypts each state block using the corresponding tweak and the key that the JFish object was
	// created with.
	Encrypt()
}

var (
	jfishMakerRegistry map[string]func(key [64]byte) JFish
)

func jfishRegister(id string, maker func(key [64]byte) JFish) {
	if jfishMakerRegistry == nil {
		jfishMakerRegistry = make(map[string]func(key [64]byte) JFish)
	}
	if _, ok := jfishMakerRegistry[id]; ok {
		panic(fmt.Sprintf("already registered a jfish for id %q", id))
	}
	jfishMakerRegistry[id] = maker
}

func MakeJFish(key [64]byte) JFish {
	if cpuid.HasExtendedFeature(cpuid.AVX2) {
		return jfishMakerRegistry["avx2"](key)
	}
	return MakeDefaultJFish(key)
}

type jFishBasic struct {
	state [8]uint64
	key   [9]uint64
	tweak [3]uint64
}

func MakeDefaultJFish(key [64]byte) JFish {
	key64 := convert.Inplace64BytesToUInt64(key[:])
	var j jFishBasic
	copy(j.key[:], key64[:])
	return &j
}

func (j *jFishBasic) NumLanes() int {
	return 1
}

func (j *jFishBasic) State(lane int) []byte {
	if lane != 0 {
		panic(fmt.Sprintf("lane %d is not available on a 1-lane JFish", lane))
	}
	return convert.InplaceUint64ToBytes(j.state[:])
}

func (j *jFishBasic) Tweak(lane int) []byte {
	if lane != 0 {
		panic(fmt.Sprintf("lane %d is not available on a 1-lane JFish", lane))
	}
	return convert.InplaceUint64ToBytes(j.tweak[0:2])
}

func (j *jFishBasic) Encrypt() {
	encrypt512(&j.state, &j.key, &j.tweak)
}

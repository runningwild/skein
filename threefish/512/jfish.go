package threefish

import (
	"fmt"

	"github.com/intel-go/cpuid"
	"github.com/runningwild/skein/types"
)

var (
	jfishMakerRegistry map[string]func() types.JFish
)

func jfishRegister(id string, maker func() types.JFish) {
	if jfishMakerRegistry == nil {
		jfishMakerRegistry = make(map[string]func() types.JFish)
	}
	if _, ok := jfishMakerRegistry[id]; ok {
		panic(fmt.Sprintf("already registered a jfish for id %q", id))
	}
	jfishMakerRegistry[id] = maker
}

func MakeJFish() types.JFish {
	if cpuid.HasExtendedFeature(cpuid.AVX2) {
		return jfishMakerRegistry["avx2"]()
	}
	return MakeDefaultJFish()
}

type jFishBasic struct {
	state [8]uint64
	key   [9]uint64
	tweak [3]uint64
}

func MakeDefaultJFish() types.JFish {
	return &jFishBasic{}
}

func (j *jFishBasic) NumLanes() int {
	return 1
}

func (j *jFishBasic) Encrypt(state, key, tweak [][]byte) {
	if len(state) != 1 || len(key) != 1 || len(tweak) != 1 {
		panic("called Encrypt without exactly one lane")
	}
	encrypt512(&j.state, &j.key, &j.tweak)
}

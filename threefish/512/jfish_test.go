package threefish_test

import (
	"testing"

	"github.com/runningwild/skein/threefish/512"
)

func TestJFish(t *testing.T) {
	var key [64]byte
	for i := range key {
		key[i] = byte(i + 100)
	}
	jf := threefish.MakeJFish(key)
	for i := 0; i < jf.NumLanes(); i++ {
		copy(jf.State(i), stateForLane(i))
		copy(jf.Tweak(i), tweakForLane(i))
	}
	jf.Encrypt()
	for i := 0; i < jf.NumLanes(); i++ {
		state := stateForLane(i)
		tbc := threefish.TweakableBlockCipher{}
		var extendedKey [72]byte
		var extendedTweak [24]byte
		copy(extendedKey[:], key[:])
		copy(extendedTweak[:], tweakForLane(i))
		tbc.Encrypt(state, extendedKey[:], extendedTweak[:])
		jfState := jf.State(i)
		if len(jfState) != len(state) {
			t.Fatalf("got state length %d, want %d", len(jf.State(i)), len(state))
		}
		for j := range jfState {
			if jfState[j] != state[j] {
				t.Errorf("state at offset %d is %d, want %d", j, jfState[j], state[j])
			}
		}
	}
}

func stateForLane(lane int) []byte {
	state := make([]byte, 64)
	for i := range state {
		state[i] = byte(i ^ lane)
	}
	return state
}

func tweakForLane(lane int) []byte {
	tweak := make([]byte, 64)
	for i := range tweak {
		tweak[i] = byte(i ^ (1 << uint(lane)))
	}
	return tweak
}

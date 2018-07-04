package threefish_test

import (
	"testing"

	"github.com/runningwild/skein/convert"
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

func BenchmarkJFish100Mb(b *testing.B) {
	b.StopTimer()
	jf := threefish.MakeJFish([64]byte{0, 1, 2})
	raw := make([]byte, 100*1024*1024/8)
	var count uint64
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		data := raw
		for len(data) > 0 {
			for i := 0; i < jf.NumLanes(); i++ {
				copy(jf.State(i), data[i*64:(i+1)*64])
				convert.InplaceBytesToUint64(jf.Tweak(i))[0] = count
				count++
			}
			jf.Encrypt()
			for i := 0; i < jf.NumLanes(); i++ {
				copy(data[i*64:(i+1)*64], jf.State(i))
			}
			data = data[jf.NumLanes()*64:]
		}
	}
}

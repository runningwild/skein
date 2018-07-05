package threefish_test

import (
	"testing"

	"github.com/runningwild/skein/convert"
	"github.com/runningwild/skein/threefish/512"
	"github.com/runningwild/skein/types"
)

func TestJFish(t *testing.T) {
	jf := threefish.MakeJFish()
	var states, keys, tweaks [][]byte
	for i := 0; i < jf.NumLanes(); i++ {
		states = append(states, stateForLane(i))
		keys = append(keys, keyForLane(i))
		tweaks = append(tweaks, tweakForLane(i))
	}
	jf.Encrypt(states, keys, tweaks)
	for i := 0; i < jf.NumLanes(); i++ {
		want := stateForLane(i)
		tbc := threefish.TweakableBlockCipher{}
		var extendedKey [72]byte
		var extendedTweak [24]byte
		copy(extendedKey[:], keys[i])
		copy(extendedTweak[:], tweakForLane(i))
		tbc.Encrypt(want, extendedKey[:], extendedTweak[:])
		for j := range states[i] {
			if states[i][j] != want[j] {
				t.Errorf("state at offset %d is %d, want %d", j, states[i][j], want[j])
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

func keyForLane(lane int) []byte {
	key := make([]byte, 64)
	return key
	for i := range key {
		key[i] = byte(i - lane)
	}
	return key
}

func tweakForLane(lane int) []byte {
	tweak := make([]byte, 64)
	for i := range tweak {
		tweak[i] = byte(i ^ (1 << uint(lane)))
	}
	return tweak
}

func BenchmarkJFish100Mb(b *testing.B) {
	for _, name := range []string{"default", "arch"} {
		b.Run(name, func(b *testing.B) {
			var jf types.JFish
			switch name {
			case "default":
				jf = threefish.MakeDefaultJFish()
			case "arch":
				jf = threefish.MakeJFish()
			default:
				b.Fatalf("unexpected test")
			}
			b.StopTimer()
			raw := make([]byte, 100*1024*1024/8)
			var count uint64
			b.StartTimer()
			states := make([][]byte, jf.NumLanes())
			keys := make([][]byte, jf.NumLanes())
			tweaks := make([][]byte, jf.NumLanes())
			tweakNums := make([][]uint64, jf.NumLanes())
			for i := 0; i < jf.NumLanes(); i++ {
				tweaks[i] = make([]byte, 16)
				tweakNums[i] = convert.InplaceBytesToUint64(tweaks[i])
			}
			for i := 0; i < b.N; i++ {
				data := raw
				for len(data) > 0 {
					for i := 0; i < jf.NumLanes(); i++ {
						states[i] = data[i*64 : (i+1)*64]
						tweakNums[i][0] = count
						count++
					}
					jf.Encrypt(states, keys, tweaks)
					data = data[jf.NumLanes()*64:]
				}
			}
		})
	}
}

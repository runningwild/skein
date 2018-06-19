package hasher_test

import (
	"fmt"
	"testing"

	skein512 "github.com/runningwild/skein/hash/512"
	"hash"
	"tmp/psha2"
	// "github.com/runningwild/skein/hash/hasher"
	// "github.com/runningwild/skein/ubi"
	// "golang.org/x/crypto/sha3"
)

type testCase struct {
	name string
	reps int
	size int
}

type hashFunction struct {
	name   string
	hasher hash.Hash
}

func BenchmarkHashes(b *testing.B) {
	for _, tc := range []testCase{
		testCase{name: "16B", reps: 1, size: 16},
		testCase{name: "1M", reps: 1, size: 1024 * 1024},
		testCase{name: "1M", reps: 10, size: 1024 * 1024},
		testCase{name: "1M", reps: 100, size: 1024 * 1024},
		testCase{name: "1M", reps: 1000, size: 1024 * 1024},
		testCase{name: "10M", reps: 1, size: 10 * 1024 * 1024},
		testCase{name: "10M", reps: 10, size: 10 * 1024 * 1024},
		testCase{name: "10M", reps: 100, size: 10 * 1024 * 1024},
		testCase{name: "100M", reps: 1, size: 100 * 1024 * 1024},
		testCase{name: "100M", reps: 10, size: 100 * 1024 * 1024},
	} {
		msg := make([]byte, tc.size)
		for _, hash := range []hashFunction{
			hashFunction{name: "skein_tree_512_1_1_2  ", hasher: skein512.NewTreeHash512(512, 11, 1, 2)},
			hashFunction{name: "psha2                 ", hasher: psha2.New()},
		} {
			b.Run(fmt.Sprintf("%s-%d*%s", hash.name, tc.reps, tc.name), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					hash.hasher.Reset()
					for rep := 0; rep < tc.reps; rep++ {
						hash.hasher.Write(msg)
					}
					hash.hasher.Sum(nil)
				}
			})
		}
	}
}

// func BenchmarkSkeinHasher_256_256_10M(b *testing.B) {
// 	b.StopTimer()
// 	u, _ := ubi.New(tf256.TweakableBlockCipher{})
// 	h := hasher.NewHasher(u, 256)
// 	msg := make([]byte, 1024)
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		for j := 0; j < 1024*10; j++ {
// 			h.Write(msg)
// 		}
// 		h.Sum(nil)
// 		h.Reset()
// 	}
// }

// func BenchmarkSkeinHasher_512_512_10M(b *testing.B) {
// 	b.StopTimer()
// 	u, _ := ubi.New(tf512.TweakableBlockCipher{})
// 	h := hasher.NewHasher(u, 512)
// 	msg := make([]byte, 1024)
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		for j := 0; j < 1024*10; j++ {
// 			h.Write(msg)
// 		}
// 		h.Sum(nil)
// 		h.Reset()
// 	}
// }

// func BenchmarkSkeinHasher_1024_1024_10M(b *testing.B) {
// 	b.StopTimer()
// 	u, _ := ubi.New(tf1024.TweakableBlockCipher{})
// 	h := hasher.NewHasher(u, 1024)
// 	msg := make([]byte, 1024)
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		for j := 0; j < 1024*10; j++ {
// 			h.Write(msg)
// 		}
// 		h.Sum(nil)
// 		h.Reset()
// 	}
// }

// func BenchmarkSHA3Hasher_256_256_10M(b *testing.B) {
// 	b.StopTimer()
// 	h := sha3.New256()
// 	msg := make([]byte, 1024)
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		for j := 0; j < 1024*10; j++ {
// 			h.Write(msg)
// 		}
// 		h.Sum(nil)
// 		h.Reset()
// 	}
// }

// func BenchmarkSHA3Hasher_512_512_10M(b *testing.B) {
// 	b.StopTimer()
// 	h := sha3.New512()
// 	msg := make([]byte, 1024)
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		for j := 0; j < 1024*10; j++ {
// 			h.Write(msg)
// 		}
// 		h.Sum(nil)
// 		h.Reset()
// 	}
// }

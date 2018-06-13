package hash_test

import (
	"testing"

	"crypto/md5"
	"crypto/sha256"
	skein1024 "github.com/runningwild/skein/hash/1024"
	skein256 "github.com/runningwild/skein/hash/256"
	skein512 "github.com/runningwild/skein/hash/512"
	"golang.org/x/crypto/sha3"
)

type testCase struct {
	name string
	size int
}

type hashFunction struct {
	name     string
	function func([]byte, int)
}

func BenchmarkHashes(b *testing.B) {
	for _, tc := range []testCase{
		testCase{name: "16B", size: 16},
		testCase{name: "1k", size: 1024},
		testCase{name: "1M", size: 1024 * 1024},
		testCase{name: "10M", size: 10 * 1024 * 1024},
	} {
		msg := make([]byte, tc.size)
		for _, hash := range []hashFunction{
			hashFunction{name: "skein_256_256  ", function: func(msg []byte, N int) {
				for i := 0; i < N; i++ {
					skein256.Hash256(256, msg)
				}
			}},
			hashFunction{name: "skein_512_512  ", function: func(msg []byte, N int) {
				for i := 0; i < N; i++ {
					skein512.Hash512(512, msg)
				}
			}},
			hashFunction{name: "skein_1024_1024", function: func(msg []byte, N int) {
				for i := 0; i < N; i++ {
					skein1024.Hash1024(1024, msg)
				}
			}},
			hashFunction{name: "md5            ", function: func(msg []byte, N int) {
				for i := 0; i < N; i++ {
					md5.Sum(msg)
				}
			}},
			hashFunction{name: "sha256         ", function: func(msg []byte, N int) {
				for i := 0; i < N; i++ {
					sha256.Sum256(msg)
				}
			}},
			hashFunction{name: "sha3_256       ", function: func(msg []byte, N int) {
				for i := 0; i < N; i++ {
					sha3.Sum256(msg)
				}
			}},
			hashFunction{name: "shake256       ", function: func(msg []byte, N int) {
				out := make([]byte, 256)
				for i := 0; i < N; i++ {
					sha3.ShakeSum256(out, msg)
				}
			}},
		} {
			b.Run(hash.name+"-"+tc.name, func(b *testing.B) {
				hash.function(msg, b.N)
			})
		}
	}
}

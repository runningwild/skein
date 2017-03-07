package hasher_test

import (
	"testing"

	"github.com/runningwild/skein/hash/hasher"
	tf1024 "github.com/runningwild/skein/threefish/1024"
	tf256 "github.com/runningwild/skein/threefish/256"
	tf512 "github.com/runningwild/skein/threefish/512"
	"github.com/runningwild/skein/ubi"
	"golang.org/x/crypto/sha3"
)

func BenchmarkSkeinHasher_256_256_10M(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(tf256.TweakableBlockCipher{})
	h := hasher.NewHasher(u, 256)
	msg := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1024*10; j++ {
			h.Write(msg)
		}
		h.Sum(nil)
		h.Reset()
	}
}

func BenchmarkSkeinHasher_512_512_10M(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(tf512.TweakableBlockCipher{})
	h := hasher.NewHasher(u, 512)
	msg := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1024*10; j++ {
			h.Write(msg)
		}
		h.Sum(nil)
		h.Reset()
	}
}

func BenchmarkSkeinHasher_1024_1024_10M(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(tf1024.TweakableBlockCipher{})
	h := hasher.NewHasher(u, 1024)
	msg := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1024*10; j++ {
			h.Write(msg)
		}
		h.Sum(nil)
		h.Reset()
	}
}

func BenchmarkSHA3Hasher_256_256_10M(b *testing.B) {
	b.StopTimer()
	h := sha3.New256()
	msg := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1024*10; j++ {
			h.Write(msg)
		}
		h.Sum(nil)
		h.Reset()
	}
}

func BenchmarkSHA3Hasher_512_512_10M(b *testing.B) {
	b.StopTimer()
	h := sha3.New512()
	msg := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1024*10; j++ {
			h.Write(msg)
		}
		h.Sum(nil)
		h.Reset()
	}
}

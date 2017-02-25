package ubi_test

import (
	"testing"

	tf1024 "github.com/runningwild/skein/threefish/1024"
	tf256 "github.com/runningwild/skein/threefish/256"
	tf512 "github.com/runningwild/skein/threefish/512"
	"github.com/runningwild/skein/ubi"
	"golang.org/x/crypto/sha3"
)

func BenchmarkSkein_256_256_16B(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(tf256.TweakableBlockCipher{})
	msg := make([]byte, 16)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u.Hash(msg, len(msg)*8, 256)
	}
}

func BenchmarkSkein_512_256_16B(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(tf512.TweakableBlockCipher{})
	msg := make([]byte, 16)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u.Hash(msg, len(msg)*8, 256)
	}
}

func BenchmarkSkein_1024_256_16B(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(tf1024.TweakableBlockCipher{})
	msg := make([]byte, 16)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u.Hash(msg, len(msg)*8, 256)
	}
}

func BenchmarkSha3_256_16B(b *testing.B) {
	b.StopTimer()
	s := sha3.New256()
	msg := make([]byte, 16)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Sum(msg)
	}
}

func BenchmarkSkein_256_256_1k(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(tf256.TweakableBlockCipher{})
	msg := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u.Hash(msg, len(msg)*8, 256)
	}
}

func BenchmarkSkein_512_256_1k(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(tf512.TweakableBlockCipher{})
	msg := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u.Hash(msg, len(msg)*8, 256)
	}
}

func BenchmarkSkein_1024_256_1k(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(tf1024.TweakableBlockCipher{})
	msg := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u.Hash(msg, len(msg)*8, 256)
	}
}

func BenchmarkSha3_256_1k(b *testing.B) {
	b.StopTimer()
	s := sha3.New256()
	msg := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Sum(msg)
	}
}

func BenchmarkSkein_256_256_1M(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(tf256.TweakableBlockCipher{})
	msg := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u.Hash(msg, len(msg)*8, 256)
	}
}

func BenchmarkSkein_512_256_1M(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(tf512.TweakableBlockCipher{})
	msg := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u.Hash(msg, len(msg)*8, 256)
	}
}

func BenchmarkSkein_1024_256_1M(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(tf1024.TweakableBlockCipher{})
	msg := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u.Hash(msg, len(msg)*8, 256)
	}
}

func BenchmarkSha3_256_1M(b *testing.B) {
	b.StopTimer()
	s := sha3.New256()
	msg := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Sum(msg)
	}
}

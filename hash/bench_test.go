package hash_test

import (
	"testing"

	skein1024 "github.com/runningwild/skein/hash/1024"
	skein256 "github.com/runningwild/skein/hash/256"
	skein512 "github.com/runningwild/skein/hash/512"
	"golang.org/x/crypto/sha3"
)

func BenchmarkSkein_256_256_16B(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 16)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		skein256.Hash256(256, msg)
	}
}

func BenchmarkSkein_512_256_16B(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 16)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		skein512.Hash512(256, msg)
	}
}

func BenchmarkSkein_1024_256_16B(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 16)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		skein1024.Hash1024(256, msg)
	}
}

func BenchmarkSha3_256_16B(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 16)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sha3.Sum256(msg)
	}
}

func BenchmarkShake_256_16B(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 16)
	out := make([]byte, 256)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sha3.ShakeSum256(msg, out)
	}
}

func BenchmarkSkein_256_256_1k(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		skein256.Hash256(256, msg)
	}
}

func BenchmarkSkein_512_256_1k(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		skein512.Hash512(256, msg)
	}
}

func BenchmarkSkein_1024_256_1k(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		skein1024.Hash1024(256, msg)
	}
}

func BenchmarkSha3_256_1k(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sha3.Sum256(msg)
	}
}

func BenchmarkShake_256_1k(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024)
	out := make([]byte, 256)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sha3.ShakeSum256(msg, out)
	}
}

func BenchmarkSkein_256_256_1M(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		skein256.Hash256(256, msg)
	}
}

func BenchmarkSkein_512_256_1M(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		skein512.Hash512(256, msg)
	}
}

func BenchmarkSkein_1024_256_1M(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		skein1024.Hash1024(256, msg)
	}
}

func BenchmarkSha3_256_1M(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sha3.Sum256(msg)
	}
}

func BenchmarkShake_256_1M(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024*1024)
	out := make([]byte, 256)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sha3.ShakeSum256(msg, out)
	}
}

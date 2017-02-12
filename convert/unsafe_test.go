package convert_test

import (
	"testing"

	"github.com/runningwild/skein/convert"
)

func BenchmarkInplace16Uint64ToBytes(b *testing.B) {
	var v [16]uint64
	for i := 0; i < b.N; i++ {
		convert.Inplace16Uint64ToBytes(v[:])
	}
}

func BenchmarkInplaceUint64ToBytes(b *testing.B) {
	v := make([]uint64, 1000)
	for i := 0; i < b.N; i++ {
		convert.InplaceUint64ToBytes(v)
	}
}

func BenchmarkInplaceBytesToUint64(b *testing.B) {
	v := make([]byte, 1000*8)
	for i := 0; i < b.N; i++ {
		convert.InplaceBytesToUint64(v)
	}
}

func BenchmarkXor8Bytes(b *testing.B) {
	x := make([]byte, 8)
	y := make([]byte, 8)
	z := make([]byte, 8)
	for i := 0; i < b.N; i++ {
		convert.Xor(x, y, z)
	}
}

func BenchmarkCipherXor64Bytes(b *testing.B) {
	x := make([]byte, 64)
	y := make([]byte, 64)
	z := make([]byte, 64)
	for i := 0; i < b.N; i++ {
		convert.XorBytes(x, y, z)
	}
}

func BenchmarkCipherXor8Bytes(b *testing.B) {
	x := make([]byte, 8)
	y := make([]byte, 8)
	z := make([]byte, 8)
	for i := 0; i < b.N; i++ {
		convert.XorBytes(x, y, z)
	}
}

func BenchmarkXor64Bytes(b *testing.B) {
	x := make([]byte, 64)
	y := make([]byte, 64)
	z := make([]byte, 64)
	for i := 0; i < b.N; i++ {
		convert.Xor(x, y, z)
	}
}

func BenchmarkNormalXor8Bytes(b *testing.B) {
	x := make([]byte, 8)
	y := make([]byte, 8)
	z := make([]byte, 8)
	for i := 0; i < b.N; i++ {
		for j := range x {
			x[j] = y[j] ^ z[j]
		}
	}
}

func BenchmarkNormalXor64Bytes(b *testing.B) {
	x := make([]byte, 64)
	y := make([]byte, 64)
	z := make([]byte, 64)
	for i := 0; i < b.N; i++ {
		for j := range x {
			x[j] = y[j] ^ z[j]
		}
	}
}

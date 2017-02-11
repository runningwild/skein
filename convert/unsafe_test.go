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

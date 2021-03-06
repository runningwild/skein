package threefish_test

import (
	"testing"

	enceve "github.com/enceve/crypto/skein/threefish"
	"github.com/runningwild/skein/threefish/1024"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncryptAndDecrypt1024(t *testing.T) {
	Convey("get the right answer for known inputs", t, func() {
		cipher := threefish.MakeCipher([128]byte{})
		state := make([]byte, 128)
		cipher.Encrypt(state, state)

		So(state, ShouldResemble, []byte{
			0xf0, 0x5c, 0x3d, 0x0a, 0x3d, 0x05, 0xb3, 0x04,
			0xf7, 0x85, 0xdd, 0xc7, 0xd1, 0xe0, 0x36, 0x01,
			0x5c, 0x8a, 0xa7, 0x6e, 0x2f, 0x21, 0x7b, 0x06,
			0xc6, 0xe1, 0x54, 0x4c, 0x0b, 0xc1, 0xa9, 0x0d,
			0xf0, 0xac, 0xcb, 0x94, 0x73, 0xc2, 0x4e, 0x0f,
			0xd5, 0x4f, 0xea, 0x68, 0x05, 0x7f, 0x43, 0x32,
			0x9c, 0xb4, 0x54, 0x76, 0x1d, 0x6d, 0xf5, 0xcf,
			0x7b, 0x2e, 0x9b, 0x36, 0x14, 0xfb, 0xd5, 0xa2,
			0x0b, 0x2e, 0x47, 0x60, 0xb4, 0x06, 0x03, 0x54,
			0x0d, 0x82, 0xea, 0xbc, 0x54, 0x82, 0xc1, 0x71,
			0xc8, 0x32, 0xaf, 0xbe, 0x68, 0x40, 0x6b, 0xc3,
			0x95, 0x00, 0x36, 0x7a, 0x59, 0x29, 0x43, 0xfa,
			0x9a, 0x5b, 0x4a, 0x43, 0x28, 0x6c, 0xa3, 0xc4,
			0xcf, 0x46, 0x10, 0x4b, 0x44, 0x31, 0x43, 0xd5,
			0x60, 0xa4, 0xb2, 0x30, 0x48, 0x83, 0x11, 0xdf,
			0x4f, 0xee, 0xf7, 0xe1, 0xdf, 0xe8, 0x39, 0x1e,
		})
		cipher.Encrypt(state, state)
		So(state, ShouldResemble, []byte{
			0x23, 0xfb, 0xe2, 0x6d, 0xf4, 0xb4, 0x84, 0x80,
			0xc8, 0xc8, 0xb2, 0xae, 0xae, 0xac, 0xfd, 0xd6,
			0xd4, 0xbc, 0x8e, 0x3f, 0x56, 0xa2, 0xde, 0x50,
			0xce, 0xfc, 0x96, 0xb2, 0x80, 0xb1, 0xd3, 0x26,
			0xe4, 0xee, 0x14, 0xd8, 0xfa, 0xa3, 0x83, 0x95,
			0xfc, 0x6a, 0xd5, 0xe5, 0x23, 0x52, 0x81, 0x09,
			0x00, 0x28, 0xd0, 0x74, 0x63, 0x12, 0x1c, 0x01,
			0x45, 0x7e, 0xef, 0x09, 0x58, 0xf9, 0xb3, 0xa8,
			0xe8, 0x1f, 0x47, 0x66, 0x32, 0x98, 0xda, 0x22,
			0x6e, 0xd4, 0xfa, 0x54, 0xc3, 0x87, 0x33, 0xfd,
			0xc4, 0xeb, 0x49, 0xba, 0xc9, 0xf5, 0x32, 0x92,
			0xd6, 0x9c, 0x45, 0x69, 0xaa, 0x58, 0xea, 0x68,
			0x10, 0xd8, 0x79, 0xb8, 0xc6, 0x73, 0xed, 0x56,
			0x85, 0xef, 0x77, 0x54, 0x53, 0x31, 0x4b, 0x16,
			0xed, 0x4f, 0x64, 0x1c, 0xbe, 0x84, 0x56, 0x22,
			0x8b, 0x9c, 0xa0, 0x00, 0xb2, 0x66, 0x2f, 0x17,
		})
		cipher.Decrypt(state, state)
		So(state, ShouldResemble, []byte{
			0xf0, 0x5c, 0x3d, 0x0a, 0x3d, 0x05, 0xb3, 0x04,
			0xf7, 0x85, 0xdd, 0xc7, 0xd1, 0xe0, 0x36, 0x01,
			0x5c, 0x8a, 0xa7, 0x6e, 0x2f, 0x21, 0x7b, 0x06,
			0xc6, 0xe1, 0x54, 0x4c, 0x0b, 0xc1, 0xa9, 0x0d,
			0xf0, 0xac, 0xcb, 0x94, 0x73, 0xc2, 0x4e, 0x0f,
			0xd5, 0x4f, 0xea, 0x68, 0x05, 0x7f, 0x43, 0x32,
			0x9c, 0xb4, 0x54, 0x76, 0x1d, 0x6d, 0xf5, 0xcf,
			0x7b, 0x2e, 0x9b, 0x36, 0x14, 0xfb, 0xd5, 0xa2,
			0x0b, 0x2e, 0x47, 0x60, 0xb4, 0x06, 0x03, 0x54,
			0x0d, 0x82, 0xea, 0xbc, 0x54, 0x82, 0xc1, 0x71,
			0xc8, 0x32, 0xaf, 0xbe, 0x68, 0x40, 0x6b, 0xc3,
			0x95, 0x00, 0x36, 0x7a, 0x59, 0x29, 0x43, 0xfa,
			0x9a, 0x5b, 0x4a, 0x43, 0x28, 0x6c, 0xa3, 0xc4,
			0xcf, 0x46, 0x10, 0x4b, 0x44, 0x31, 0x43, 0xd5,
			0x60, 0xa4, 0xb2, 0x30, 0x48, 0x83, 0x11, 0xdf,
			0x4f, 0xee, 0xf7, 0xe1, 0xdf, 0xe8, 0x39, 0x1e,
		})
		cipher.Decrypt(state, state)
		So(state, ShouldResemble, []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		})
	})
}

func BenchmarkEncrypt1024Block(b *testing.B) {
	var state [128]byte
	var key [136]byte
	var tweak [24]byte
	for i := 0; i < b.N; i++ {
		threefish.Encrypt(state[:], key[:], tweak[:])
	}
}

func BenchmarkEncrypt1024Block_enceve(b *testing.B) {
	var state [16]uint64
	var key [17]uint64
	var tweak [3]uint64
	for i := 0; i < b.N; i++ {
		enceve.Encrypt1024(&state, &key, &tweak)
	}
}

func BenchmarkDecrypt1024Block(b *testing.B) {
	var state [128]byte
	var key [136]byte
	var tweak [24]byte
	for i := 0; i < b.N; i++ {
		threefish.Decrypt(state[:], key[:], tweak[:])
	}
}

func BenchmarkDecrypt1024Block_enceve(b *testing.B) {
	var state [16]uint64
	var key [17]uint64
	var tweak [3]uint64
	for i := 0; i < b.N; i++ {
		enceve.Decrypt1024(&state, &key, &tweak)
	}
}

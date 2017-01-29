package threefish_test

import (
	"testing"

	enceve "github.com/enceve/crypto/skein/threefish"
	"github.com/runningwild/skein/threefish"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncryptAndDecrypt256(t *testing.T) {
	Convey("get the right answer for known inputs", t, func() {
		cipher := threefish.MakeCipher256(make([]byte, 32))
		state := make([]byte, 32)
		cipher.Encrypt(state, state)
		So(state, ShouldResemble, []byte{
			0x84, 0xda, 0x2a, 0x1f, 0x8b, 0xea, 0xee, 0x94,
			0x70, 0x66, 0xae, 0x3e, 0x31, 0x03, 0xf1, 0xad,
			0x53, 0x6d, 0xb1, 0xf4, 0xa1, 0x19, 0x24, 0x95,
			0x11, 0x6b, 0x9f, 0x3c, 0xe6, 0x13, 0x3f, 0xd8,
		})
		cipher.Encrypt(state, state)
		So(state, ShouldResemble, []byte{
			0x43, 0x5f, 0xdc, 0xf2, 0xfd, 0x3a, 0xb9, 0x35,
			0x71, 0x3f, 0x12, 0x6b, 0xfb, 0x32, 0x20, 0x3b,
			0x56, 0x2b, 0xf2, 0xd3, 0x1f, 0x26, 0x31, 0x46,
			0xaf, 0xa5, 0x34, 0x60, 0x3f, 0x63, 0x97, 0x20,
		})
		cipher.Decrypt(state, state)
		So(state, ShouldResemble, []byte{
			0x84, 0xda, 0x2a, 0x1f, 0x8b, 0xea, 0xee, 0x94,
			0x70, 0x66, 0xae, 0x3e, 0x31, 0x03, 0xf1, 0xad,
			0x53, 0x6d, 0xb1, 0xf4, 0xa1, 0x19, 0x24, 0x95,
			0x11, 0x6b, 0x9f, 0x3c, 0xe6, 0x13, 0x3f, 0xd8,
		})
		cipher.Decrypt(state, state)
		So(state, ShouldResemble, []byte{
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		})
	})
}

func BenchmarkEncrypt256Block(b *testing.B) {
	cipher := threefish.MakeCipher256(make([]byte, 32))
	block := make([]byte, 32)
	for i := 0; i < b.N; i++ {
		cipher.Encrypt(block, block)
	}
}

func BenchmarkEncrypt256Block_enceve(b *testing.B) {
	var state [4]uint64
	var key [5]uint64
	var tweak [3]uint64
	for i := 0; i < b.N; i++ {
		enceve.Encrypt256(&state, &key, &tweak)
	}
}

func BenchmarkDecrypt256Block(b *testing.B) {
	cipher := threefish.MakeCipher256(make([]byte, 32))
	block := make([]byte, 32)
	for i := 0; i < b.N; i++ {
		cipher.Decrypt(block, block)
	}
}

func BenchmarkDecrypt256Block_enceve(b *testing.B) {
	var state [4]uint64
	var key [5]uint64
	var tweak [3]uint64
	for i := 0; i < b.N; i++ {
		enceve.Decrypt256(&state, &key, &tweak)
	}
}

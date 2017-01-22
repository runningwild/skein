package skein

import (
	"testing"

	enceve "github.com/enceve/crypto/skein/threefish"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncryptAndDecrypt(t *testing.T) {
	Convey("get the right answer for known inputs", t, func() {
		var b block256
		b.encrypt()
		So(b.state, ShouldResemble, [4]uint64{0x94eeea8b1f2ada84, 0xadf103313eae6670, 0x952419a1f4b16d53, 0xd83f13e63c9f6b11})
		b.encrypt()
		So(b.state, ShouldResemble, [4]uint64{0x35b93afdf2dc5f43, 0x3b2032fb6b123f71, 0x4631261fd3f22b56, 0x2097633f6034a5af})
		b.decrypt()
		So(b.state, ShouldResemble, [4]uint64{0x94eeea8b1f2ada84, 0xadf103313eae6670, 0x952419a1f4b16d53, 0xd83f13e63c9f6b11})
		b.decrypt()
		So(b.state, ShouldResemble, [4]uint64{0, 0, 0, 0})
	})
}

func BenchmarkEncryptBlock(b *testing.B) {
	var block block256
	for i := 0; i < b.N; i++ {
		block.encrypt()
	}
}

func BenchmarkEncryptBlock_enceve(b *testing.B) {
	var block [4]uint64
	var keys [5]uint64
	var tweak [3]uint64
	for i := 0; i < b.N; i++ {
		enceve.Encrypt256(&block, &keys, &tweak)
	}
}

func BenchmarkEncryptBlockSlow(b *testing.B) {
	var block block256_slow
	for i := 0; i < b.N; i++ {
		block.encrypt()
	}
}

func BenchmarkDecryptBlock(b *testing.B) {
	var block block256
	for i := 0; i < b.N; i++ {
		block.decrypt()
	}
}

func BenchmarkDecryptBlockSlow(b *testing.B) {
	var block block256_slow
	for i := 0; i < b.N; i++ {
		block.decrypt()
	}
}

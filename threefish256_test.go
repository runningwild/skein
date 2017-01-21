package skein

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMix(t *testing.T) {
	Convey("rounds work", t, func() {
		var stateGot, stateWant [4]uint64
		stateGot[0], stateGot[1], stateGot[2], stateGot[3] = 4, 5, 6, 7
		stateWant[0], stateWant[1], stateWant[2], stateWant[3] = 4, 5, 6, 7

		fmt.Printf("%v %v\n", stateGot, stateWant)
		fourRoundsA(&stateGot)
		fmt.Printf("%v %v\n", stateGot, stateWant)
		fourRoundsA_slow(&stateWant)
		fmt.Printf("%v %v\n", stateGot, stateWant)
		So(stateGot, ShouldEqual, stateWant)
		fourRoundsA(&stateGot)
		fourRoundsA_slow(&stateWant)
		So(stateGot, ShouldEqual, stateWant)
		fourRoundsA(&stateGot)
		fourRoundsA_slow(&stateWant)
		So(stateGot, ShouldEqual, stateWant)
	})
}

func mix256_ref(rot uint, x0, x1 uint64) (y0, y1 uint64) {
	y0 = x0 + x1
	y1 = (x1 << rot) | (x1 >> (64 - rot))
	y1 = y1 ^ y0
	return
}

func TestEncryptAndDecrypt(t *testing.T) {
	Convey("get the right answer for known inputs", t, func() {
		var b block256
		b.encrypt()
		So(b.state, ShouldEqual, [4]uint64{0x94eeea8b1f2ada84, 0xadf103313eae6670, 0x952419a1f4b16d53, 0xd83f13e63c9f6b11})
		b.encrypt()
		So(b.state, ShouldEqual, [4]uint64{0x35b93afdf2dc5f43, 0x3b2032fb6b123f71, 0x4631261fd3f22b56, 0x2097633f6034a5af})
		b.decrypt()
		So(b.state, ShouldEqual, [4]uint64{0x94eeea8b1f2ada84, 0xadf103313eae6670, 0x952419a1f4b16d53, 0xd83f13e63c9f6b11})
		b.decrypt()
		So(b.state, ShouldEqual, [4]uint64{0, 0, 0, 0})
	})
}

func BenchmarkEncryptBlock(b *testing.B) {
	var block block256
	for i := 0; i < b.N; i++ {
		block.encrypt()
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

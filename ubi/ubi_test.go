package ubi_test

import (
	"testing"

	"github.com/runningwild/skein/threefish/256"
	"github.com/runningwild/skein/ubi"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSkeinGeneral(t *testing.T) {
	sampleMessages := [][]byte{
		[]byte(""),
		[]byte("Buttons is an adorable dog"),
		[]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."),
	}
	giantMess := make([]byte, 1024*1024)
	for i := range giantMess {
		giantMess[i] = byte(i)
	}
	sampleMessages = append(sampleMessages, giantMess)

	Convey("ubi.Skein", t, func() {
		u, err := ubi.New(threefish.Encrypt, 256)
		So(err, ShouldBeNil)

		Convey("doesn't corrupt input message", func() {
			msg := make([]byte, 100)
			u.Skein(msg, 256)
			for i := range msg {
				So(msg[i], ShouldBeZeroValue)
			}
		})

		Convey("matches skein256_N", func() {
			for _, msg := range sampleMessages {
				for _, outlen := range []uint64{10, 100, 1000, 8000, 10 * 1024} {
					h := u.Skein(msg, outlen)
					So(h[:], ShouldResemble, ubi.Skein256_N(msg, outlen))
				}
			}
		})

		Convey("passes test vectors", func() {
			So(u.Skein(make([]byte, 128), 512), ShouldResemble, []byte{
				0x91, 0x8A, 0x1B, 0x6D, 0x20, 0x01, 0x5D, 0x0B, 0xF5, 0x3C, 0xF4, 0xFD, 0xD3, 0x9E, 0x28, 0xD8,
				0xBC, 0x55, 0x04, 0xA9, 0x6C, 0x1D, 0x31, 0x0A, 0xD5, 0xAD, 0xB1, 0x5D, 0xCD, 0xDE, 0xA2, 0x70,
				0x18, 0x4F, 0x94, 0x67, 0x45, 0x1C, 0xD9, 0x7B, 0xC6, 0x24, 0xD3, 0x08, 0x83, 0xA0, 0x06, 0x33,
				0x64, 0x57, 0x81, 0x5A, 0x88, 0xA9, 0xFE, 0xB4, 0x49, 0x46, 0x0E, 0x4B, 0x42, 0xD9, 0x66, 0xAC,
			})
		})
	})
}

func BenchmarkSkeinGeneral_256_256_16B(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(threefish.Encrypt, 256)
	msg := make([]byte, 16)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u.Skein(msg, 256)
	}
}

func BenchmarkSkeinGeneral_256_256_1M(b *testing.B) {
	b.StopTimer()
	u, _ := ubi.New(threefish.Encrypt, 256)
	msg := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u.Skein(msg, 256)
	}
}

package skein

import (
	"testing"

	"crypto/sha1"
	"crypto/sha256"

	"golang.org/x/crypto/sha3"

	enceve "github.com/enceve/crypto/skein/skein256"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUBI(t *testing.T) {
	// Convey("ubi", t, func() {
	// 	out := ubi256([4]uint64{}, []byte{
	// 		0x53, 0x48, 0x41, 0x33, // SHA1
	// 		0x01, 0x00, // Version Number
	// 		0x00, 0x00, // Reserved
	// 		0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Output size in bits (128)
	// 		0x00, 0x00, 0x00, // Tree params
	// 		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Reserved,
	// 	}, [2]uint64{0, 4 << (120 - 64)})
	// 	So(fmt.Sprintf("%x", out), ShouldEqual, "x")
	// })
}

func TestSkein(t *testing.T) {
	Convey("skein256_256 passes test vectors", t, func() {
		So(skein256_256([]byte{}), ShouldEqual, [32]byte{
			0xC8, 0x87, 0x70, 0x87, 0xDA, 0x56, 0xE0, 0x72, 0x87, 0x0D, 0xAA, 0x84, 0x3F, 0x17, 0x6E, 0x94,
			0x53, 0x11, 0x59, 0x29, 0x09, 0x4C, 0x3A, 0x40, 0xC4, 0x63, 0xA1, 0x96, 0xC2, 0x9B, 0xF7, 0xBA,
		})

		So(skein256_256([]byte{
			0x00,
		}), ShouldEqual, [32]byte{
			0x34, 0xE2, 0xB6, 0x5B, 0xF0, 0xBE, 0x66, 0x7C, 0xA5, 0xDE, 0xBA, 0x82, 0xC3, 0x7C, 0xB2, 0x53,
			0xEB, 0x9F, 0x84, 0x74, 0xF3, 0x42, 0x6B, 0xA6, 0x22, 0xA2, 0x52, 0x19, 0xFD, 0x18, 0x24, 0x33,
		})

		So(skein256_256([]byte{
			0xFF,
		}), ShouldEqual, [32]byte{
			0X0B, 0X98, 0XDC, 0XD1, 0X98, 0XEA, 0X0E, 0X50, 0XA7, 0XA2, 0X44, 0XC4, 0X44, 0XE2, 0X5C, 0X23,
			0XDA, 0X30, 0XC1, 0X0F, 0XC9, 0XA1, 0XF2, 0X70, 0XA6, 0X63, 0X7F, 0X1F, 0X34, 0XE6, 0X7E, 0XD2,
		})

		So(skein256_256([]byte{
			0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1, 0xF0,
			0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8, 0xE7, 0xE6, 0xE5, 0xE4, 0xE3, 0xE2, 0xE1, 0xE0,
		}), ShouldEqual, [32]byte{
			0x8D, 0x0F, 0xA4, 0xEF, 0x77, 0x7F, 0xD7, 0x59, 0xDF, 0xD4, 0x04, 0x4E, 0x6F, 0x6A, 0x5A, 0xC3,
			0xC7, 0x74, 0xAE, 0xC9, 0x43, 0xDC, 0xFC, 0x07, 0x92, 0x7B, 0x72, 0x3B, 0x5D, 0xBF, 0x40, 0x8B,
		})

		So(skein256_256([]byte{
			0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1, 0xF0,
			0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8, 0xE7, 0xE6, 0xE5, 0xE4, 0xE3, 0xE2, 0xE1, 0xE0,
			0xDF, 0xDE, 0xDD, 0xDC, 0xDB, 0xDA, 0xD9, 0xD8, 0xD7, 0xD6, 0xD5, 0xD4, 0xD3, 0xD2, 0xD1, 0xD0,
			0xCF, 0xCE, 0xCD, 0xCC, 0xCB, 0xCA, 0xC9, 0xC8, 0xC7, 0xC6, 0xC5, 0xC4, 0xC3, 0xC2, 0xC1, 0xC0,
		}), ShouldEqual, [32]byte{
			0xDF, 0x28, 0xE9, 0x16, 0x63, 0x0D, 0x0B, 0x44, 0xC4, 0xA8, 0x49, 0xDC, 0x9A, 0x02, 0xF0, 0x7A,
			0x07, 0xCB, 0x30, 0xF7, 0x32, 0x31, 0x82, 0x56, 0xB1, 0x5D, 0x86, 0x5A, 0xC4, 0xAE, 0x16, 0x2F,
		})
	})

	Convey("skein256_N passes test vectors", t, func() {
		So(skein256_N([]byte{}, 256), ShouldResemble, []byte{
			0xC8, 0x87, 0x70, 0x87, 0xDA, 0x56, 0xE0, 0x72, 0x87, 0x0D, 0xAA, 0x84, 0x3F, 0x17, 0x6E, 0x94,
			0x53, 0x11, 0x59, 0x29, 0x09, 0x4C, 0x3A, 0x40, 0xC4, 0x63, 0xA1, 0x96, 0xC2, 0x9B, 0xF7, 0xBA,
		})

		So(skein256_N([]byte{
			0x00,
		}, 256), ShouldResemble, []byte{
			0x34, 0xE2, 0xB6, 0x5B, 0xF0, 0xBE, 0x66, 0x7C, 0xA5, 0xDE, 0xBA, 0x82, 0xC3, 0x7C, 0xB2, 0x53,
			0xEB, 0x9F, 0x84, 0x74, 0xF3, 0x42, 0x6B, 0xA6, 0x22, 0xA2, 0x52, 0x19, 0xFD, 0x18, 0x24, 0x33,
		})

		So(skein256_N([]byte{
			0xFF,
		}, 256), ShouldResemble, []byte{
			0X0B, 0X98, 0XDC, 0XD1, 0X98, 0XEA, 0X0E, 0X50, 0XA7, 0XA2, 0X44, 0XC4, 0X44, 0XE2, 0X5C, 0X23,
			0XDA, 0X30, 0XC1, 0X0F, 0XC9, 0XA1, 0XF2, 0X70, 0XA6, 0X63, 0X7F, 0X1F, 0X34, 0XE6, 0X7E, 0XD2,
		})

		So(skein256_N([]byte{
			0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1, 0xF0,
			0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8, 0xE7, 0xE6, 0xE5, 0xE4, 0xE3, 0xE2, 0xE1, 0xE0,
		}, 256), ShouldResemble, []byte{
			0x8D, 0x0F, 0xA4, 0xEF, 0x77, 0x7F, 0xD7, 0x59, 0xDF, 0xD4, 0x04, 0x4E, 0x6F, 0x6A, 0x5A, 0xC3,
			0xC7, 0x74, 0xAE, 0xC9, 0x43, 0xDC, 0xFC, 0x07, 0x92, 0x7B, 0x72, 0x3B, 0x5D, 0xBF, 0x40, 0x8B,
		})

		So(skein256_N([]byte{
			0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1, 0xF0,
			0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8, 0xE7, 0xE6, 0xE5, 0xE4, 0xE3, 0xE2, 0xE1, 0xE0,
			0xDF, 0xDE, 0xDD, 0xDC, 0xDB, 0xDA, 0xD9, 0xD8, 0xD7, 0xD6, 0xD5, 0xD4, 0xD3, 0xD2, 0xD1, 0xD0,
			0xCF, 0xCE, 0xCD, 0xCC, 0xCB, 0xCA, 0xC9, 0xC8, 0xC7, 0xC6, 0xC5, 0xC4, 0xC3, 0xC2, 0xC1, 0xC0,
		}, 256), ShouldResemble, []byte{
			0xDF, 0x28, 0xE9, 0x16, 0x63, 0x0D, 0x0B, 0x44, 0xC4, 0xA8, 0x49, 0xDC, 0x9A, 0x02, 0xF0, 0x7A,
			0x07, 0xCB, 0x30, 0xF7, 0x32, 0x31, 0x82, 0x56, 0xB1, 0x5D, 0x86, 0x5A, 0xC4, 0xAE, 0x16, 0x2F,
		})

		So(skein256_N([]byte{
			0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8, 0xF7, 0xF6, 0xF5, 0xF4, 0xF3, 0xF2, 0xF1, 0xF0,
			0xEF, 0xEE, 0xED, 0xEC, 0xEB, 0xEA, 0xE9, 0xE8, 0xE7, 0xE6, 0xE5, 0xE4, 0xE3, 0xE2, 0xE1, 0xE0,
			0xDF, 0xDE, 0xDD, 0xDC, 0xDB, 0xDA, 0xD9, 0xD8, 0xD7, 0xD6, 0xD5, 0xD4, 0xD3, 0xD2, 0xD1, 0xD0,
			0xCF, 0xCE, 0xCD, 0xCC, 0xCB, 0xCA, 0xC9, 0xC8, 0xC7, 0xC6, 0xC5, 0xC4, 0xC3, 0xC2, 0xC1, 0xC0,
			0xBF, 0xBE, 0xBD, 0xBC, 0xBB, 0xBA, 0xB9, 0xB8, 0xB7, 0xB6, 0xB5, 0xB4, 0xB3, 0xB2, 0xB1, 0xB0,
			0xAF, 0xAE, 0xAD, 0xAC, 0xAB, 0xAA, 0xA9, 0xA8, 0xA7, 0xA6, 0xA5, 0xA4, 0xA3, 0xA2, 0xA1, 0xA0,
			0x9F, 0x9E, 0x9D, 0x9C, 0x9B, 0x9A, 0x99, 0x98, 0x97, 0x96, 0x95, 0x94, 0x93, 0x92, 0x91, 0x90,
			0x8F, 0x8E, 0x8D, 0x8C, 0x8B, 0x8A, 0x89, 0x88, 0x87, 0x86, 0x85, 0x84, 0x83, 0x82, 0x81, 0x80,
		}, 2056), ShouldResemble, []byte{
			0xB6, 0x18, 0x96, 0xA9, 0xC8, 0xE9, 0x39, 0xF3, 0x0B, 0x55, 0x48, 0x11, 0x22, 0x17, 0x32, 0x17,
			0x49, 0xA1, 0xFE, 0x88, 0xF3, 0xB1, 0x81, 0x48, 0x97, 0xBC, 0x5D, 0x47, 0x09, 0x85, 0xFF, 0x50,
			0xB6, 0xF6, 0x6B, 0xCF, 0x73, 0xAD, 0x68, 0x2B, 0x49, 0x7A, 0xE7, 0x37, 0x33, 0x18, 0x6A, 0xEA,
			0xBF, 0xC1, 0x93, 0x8F, 0x9B, 0x43, 0x8A, 0xF4, 0x32, 0xEF, 0x91, 0xDE, 0x99, 0x48, 0x3F, 0x3A,
			0xB0, 0xDC, 0x71, 0xB3, 0xC9, 0x09, 0x31, 0x1A, 0x47, 0xB6, 0xDC, 0x77, 0xD9, 0x74, 0xDB, 0x3A,
			0xF6, 0x83, 0x47, 0x0C, 0x88, 0x43, 0x89, 0xDB, 0x6F, 0xA3, 0xC2, 0x19, 0xD9, 0xAD, 0x82, 0x2B,
			0x09, 0xE2, 0x44, 0x48, 0x12, 0x52, 0xE4, 0xDA, 0xA9, 0xCC, 0xCE, 0x6A, 0x3F, 0x70, 0x9C, 0x57,
			0xF3, 0x55, 0x3C, 0x6D, 0x94, 0x1F, 0x70, 0x26, 0x57, 0x0A, 0xD5, 0x74, 0xE5, 0x82, 0x99, 0xB1,
			0x92, 0x78, 0x8C, 0xA8, 0x76, 0x05, 0xD4, 0x42, 0x94, 0x29, 0x90, 0xCC, 0x88, 0x35, 0xAD, 0x89,
			0xDE, 0xA3, 0x50, 0x47, 0x95, 0x34, 0x4B, 0xB2, 0x38, 0xA9, 0x11, 0x35, 0x69, 0x1A, 0xA7, 0x70,
			0xAE, 0x6A, 0xDF, 0xE9, 0xBC, 0xFC, 0xD6, 0xC2, 0x8C, 0x78, 0x0A, 0x47, 0xB7, 0x8D, 0x24, 0x85,
			0x65, 0xCE, 0x49, 0xF6, 0xC2, 0xDF, 0xEC, 0xB8, 0xD8, 0xC2, 0xA4, 0xE8, 0x95, 0x84, 0x5B, 0xED,
			0x0D, 0xC7, 0x37, 0x46, 0xBD, 0x56, 0xA6, 0x39, 0x56, 0x55, 0x68, 0x1C, 0xDE, 0x63, 0x68, 0xAA,
			0x4C, 0x58, 0x03, 0x29, 0xBA, 0xF2, 0x84, 0xEF, 0xB3, 0x10, 0x4E, 0xD4, 0x98, 0x4E, 0x4B, 0x07,
			0xEA, 0x8E, 0x9C, 0x78, 0xFB, 0x99, 0x49, 0x29, 0x08, 0x44, 0x4F, 0x4F, 0x78, 0x16, 0x02, 0xC3,
			0x79, 0x2A, 0x9B, 0xB1, 0x51, 0x26, 0xC2, 0xDA, 0x28, 0xD9, 0x39, 0x0E, 0xE5, 0x6F, 0xF0, 0xFC,
			0x22,
		})
	})
}

func BenchmarkSkein256_256_16B(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 16)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		skein256_256(msg)
	}
}

func BenchmarkSkein_enceve_256_256_16B(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 16)
	var out [32]byte
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		enceve.Sum256(&out, msg, nil)
	}
}

func BenchmarkSkein256_256_1M(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		skein256_256(msg)
	}
}

func BenchmarkSkein_enceve_256_256_1M(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024*1024)
	var out [32]byte
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		enceve.Sum256(&out, msg, nil)
	}
}

func BenchmarkSHA1_1M(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sha1.Sum(msg)
	}
}

func BenchmarkSHA256_1M(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sha256.Sum256(msg)
	}
}

func BenchmarkSHA3_1M(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sha3.Sum256(msg)
	}
}

func BenchmarkSkein256_256_1G(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024*1024*1024)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		skein256_256(msg)
	}
}

func BenchmarkSkein_enceve_256_256_1G(b *testing.B) {
	b.StopTimer()
	msg := make([]byte, 1024*1024*1024)
	var out [32]byte
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		enceve.Sum256(&out, msg, nil)
	}
}

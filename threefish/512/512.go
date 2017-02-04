package threefish

import (
	"github.com/runningwild/skein/convert"
)

const c240 = 0x1bd11bdaa9fc1a22

// Cipher implements the go standard cipher.Block interface using the threefish 512-bit cipher.
type Cipher struct {
	key   [9]uint64 // The user-defined key is the first 8 values, we add the 9th.
	tweak [3]uint64 // The user-defined tweak is the first 2 values, we add the 3rd.
}

// MakeCipher returns a Cipher object using the specified key.
func MakeCipher(key [64]byte) *Cipher {
	var c Cipher
	key64 := convert.Inplace64BytesToUInt64(key[:])
	copy(c.key[:], key64[:])
	return &c
}

func (c *Cipher) BlockSize() int { return 64 }

func (c *Cipher) Encrypt(dst, src []byte) {
	if len(src) != len(dst) || len(src) != 64 {
		panic("Cipher.Encrypt requires src and dst slices must both be exactly 64 bytes")
	}
	if &dst[0] != &src[0] {
		copy(dst, src)
	}
	encrypt512(convert.Inplace64BytesToUInt64(dst), &c.key, &c.tweak)
}

func (c *Cipher) Decrypt(dst, src []byte) {
	if len(src) != len(dst) || len(src) != 64 {
		panic("Cipher.Decrypt requires src and dst slices must both be exactly 64 bytes")
	}
	if &dst[0] != &src[0] {
		copy(dst, src)
	}
	decrypt512(convert.Inplace64BytesToUInt64(dst), &c.key, &c.tweak)
}

// Encrypt encrypts a single block of 64 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 4 elements of key, the fifth value is used internally.
// Similarly, the tweak should be contained in the first two elements of tweak, the third value is
// used internally.
func Encrypt(data []byte, key *[9]uint64, tweak *[3]uint64) {
	if len(data) != 64 {
		panic("Encrypt requires that data is exactly 64 bytes")
	}
	encrypt512(convert.Inplace64BytesToUInt64(data), key, tweak)
}

// Decrypt decrypts a single block of 64 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 4 elements of key, the fifth value is used internally.
// Similarly, the tweak should be contained in the first two elements of tweak, the third value is
// used internally.
func Decrypt(data []byte, key *[9]uint64, tweak *[3]uint64) {
	if len(data) != 64 {
		panic("Decrypt requires that data is exactly 64 bytes")
	}
	decrypt512(convert.Inplace64BytesToUInt64(data), key, tweak)
}
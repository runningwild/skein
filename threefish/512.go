package threefish

import (
	"github.com/runningwild/skein/convert"
)

// Cipher512 implements the go standard cipher.Block interface using the threefish 512-bit cipher.
type Cipher512 struct {
	key   [9]uint64 // The user-defined key is the first 8 values, we add the 9th.
	tweak [3]uint64 // The user-defined tweak is the first 2 values, we add the 3rd.
}

// MakeCipher512 returns a Cipher512 object using the specified key.  key must be exactly 64 bytes
// or this function will panic.
func MakeCipher512(key []byte) *Cipher512 {
	if len(key) != 64 {
		panic("MakeCipher512 can only be called with a key of exactly 64 bytes")
	}
	var c Cipher512
	key64 := convert.Inplace64BytesToUInt64(key)
	copy(c.key[:], key64[:])
	return &c
}

func (c *Cipher512) BlockSize() int { return 64 }

func (c *Cipher512) Encrypt(dst, src []byte) {
	if len(src) != len(dst) || len(src) != 64 {
		panic("Cipher512.Encrypt requires src and dst slices must both be exactly 64 bytes")
	}
	if &dst[0] != &src[0] {
		copy(dst, src)
	}
	Encrypt512(dst, &c.key, &c.tweak)
}

func (c *Cipher512) Decrypt(dst, src []byte) {
	if len(src) != len(dst) || len(src) != 64 {
		panic("Cipher512.Decrypt requires src and dst slices must both be exactly 64 bytes")
	}
	if &dst[0] != &src[0] {
		copy(dst, src)
	}
	decrypt512Simple(convert.Inplace64BytesToUInt64(dst), &c.key, &c.tweak)
}

// Encrypt512 encrypts a single block of 64 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 4 elements of key, the fifth value is used internally.
// Similarly, the tweak should be contained in the first two elements of tweak, the third value is
// used internally.
func Encrypt512(data []byte, key *[9]uint64, tweak *[3]uint64) {
	if len(data) != 64 {
		panic("Encrypt512 requires that data is exactly 64 bytes")
	}
	encrypt512(convert.Inplace64BytesToUInt64(data), key, tweak)
}

// Decrypt512 decrypts a single block of 64 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 4 elements of key, the fifth value is used internally.
// Similarly, the tweak should be contained in the first two elements of tweak, the third value is
// used internally.
func Decrypt512(data []byte, key *[9]uint64, tweak *[3]uint64) {
	if len(data) != 64 {
		panic("Decrypt512 requires that data is exactly 64 bytes")
	}
	decrypt512(convert.Inplace64BytesToUInt64(data), key, tweak)
}

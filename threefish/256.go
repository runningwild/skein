package threefish

import (
	"github.com/runningwild/skein/convert"
)

// Cipher256 implements the go standard cipher.Block interface using the threefish 256-bit cipher.
type Cipher256 struct {
	key   [5]uint64 // The user-defined key is the first 4 values, we add the 5th.
	tweak [3]uint64 // The user-defined tweak is the first 2 values, we add the 3rd.
}

// MakeCipher256 returns a Cipher256 object using the specified key.  key must be exactly 32 bytes
// or this function will panic.
func MakeCipher256(key []byte) *Cipher256 {
	if len(key) != 32 {
		panic("MakeCipher256 can only be called with a key of exactly 32 bytes")
	}
	var c Cipher256
	key64 := convert.Inplace32BytesToUInt64(key)
	copy(c.key[:], key64[:])
	return &c
}

func (c *Cipher256) BlockSize() int { return 32 }

func (c *Cipher256) Encrypt(dst, src []byte) {
	if len(src) != len(dst) || len(src) != 32 {
		panic("Cipher256.Encrypt requires src and dst slices must both be exactly 32 bytes")
	}
	if &dst[0] != &src[0] {
		copy(dst, src)
	}
	Encrypt256(dst, &c.key, &c.tweak)
}

func (c *Cipher256) Decrypt(dst, src []byte) {
	if len(src) != len(dst) || len(src) != 32 {
		panic("Cipher256.Decrypt requires src and dst slices must both be exactly 32 bytes")
	}
	if &dst[0] != &src[0] {
		copy(dst, src)
	}
	decrypt256Simple(convert.Inplace32BytesToUInt64(dst), &c.key, &c.tweak)
}

// Encrypt256 encrypts a single block of 32 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 4 elements of key, the fifth value is used internally.
// Similarly, the tweak should be contained in the first two elements of tweak, the third value is
// used internally.
func Encrypt256(data []byte, key *[5]uint64, tweak *[3]uint64) {
	if len(data) != 32 {
		panic("Encrypt256 requires that data is exactly 32 bytes")
	}
	encrypt256(convert.Inplace32BytesToUInt64(data), key, tweak)
}

// Decrypt256 decrypts a single block of 32 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 4 elements of key, the fifth value is used internally.
// Similarly, the tweak should be contained in the first two elements of tweak, the third value is
// used internally.
func Decrypt256(data []byte, key *[5]uint64, tweak *[3]uint64) {
	if len(data) != 32 {
		panic("Decrypt256 requires that data is exactly 32 bytes")
	}
	decrypt256(convert.Inplace32BytesToUInt64(data), key, tweak)
}

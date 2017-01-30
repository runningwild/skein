package threefish

import (
	"github.com/runningwild/skein/convert"
)

// Cipher1024 implements the go standard cipher.Block interface using the threefish 1024-bit cipher.
type Cipher1024 struct {
	key   [17]uint64 // The user-defined key is the first 16 values, we add the 17th.
	tweak [3]uint64  // The user-defined tweak is the first 2 values, we add the 3rd.
}

// MakeCipher1024 returns a Cipher1024 object using the specified key.  key must be exactly 128
// bytes or this function will panic.
func MakeCipher1024(key []byte) *Cipher1024 {
	if len(key) != 128 {
		panic("MakeCipher1024 can only be called with a key of exactly 128 bytes")
	}
	var c Cipher1024
	key128 := convert.Inplace128BytesToUInt64(key)
	copy(c.key[:], key128[:])
	return &c
}

func (c *Cipher1024) BlockSize() int { return 128 }

func (c *Cipher1024) Encrypt(dst, src []byte) {
	if len(src) != len(dst) || len(src) != 128 {
		panic("Cipher1024.Encrypt requires src and dst slices must both be exactly 128 bytes")
	}
	if &dst[0] != &src[0] {
		copy(dst, src)
	}
	Encrypt1024(dst, &c.key, &c.tweak)
}

func (c *Cipher1024) Decrypt(dst, src []byte) {
	if len(src) != len(dst) || len(src) != 128 {
		panic("Cipher1024.Decrypt requires src and dst slices must both be exactly 128 bytes")
	}
	if &dst[0] != &src[0] {
		copy(dst, src)
	}
	decrypt1024Simple(convert.Inplace128BytesToUInt64(dst), &c.key, &c.tweak)
}

// Encrypt1024 encrypts a single block of 128 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 4 elements of key, the fifth value is used internally.
// Similarly, the tweak should be contained in the first two elements of tweak, the third value is
// used internally.
func Encrypt1024(data []byte, key *[17]uint64, tweak *[3]uint64) {
	if len(data) != 128 {
		panic("Encrypt1024 requires that data is exactly 128 bytes")
	}
	encrypt1024(convert.Inplace128BytesToUInt64(data), key, tweak)
}

// Decrypt1024 decrypts a single block of 128 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 4 elements of key, the fifth value is used internally.
// Similarly, the tweak should be contained in the first two elements of tweak, the third value is
// used internally.
func Decrypt1024(data []byte, key *[17]uint64, tweak *[3]uint64) {
	if len(data) != 128 {
		panic("Decrypt1024 requires that data is exactly 128 bytes")
	}
	decrypt1024(convert.Inplace128BytesToUInt64(data), key, tweak)
}

package threefish

import (
	"github.com/runningwild/skein/convert"
)

const c240 = 0x1bd11bdaa9fc1a22

// Cipher implements the go standard cipher.Block interface using the threefish 256-bit cipher.
type Cipher struct {
	key   [5]uint64 // The user-defined key is the first 4 values, we add the 5th.
	tweak [3]uint64 // The user-defined tweak is the first 2 values, we add the 3rd.
}

// MakeCipher returns a Cipher object using the specified key.
func MakeCipher(key [32]byte) *Cipher {
	var c Cipher
	key64 := convert.Inplace32BytesToUInt64(key[:])
	copy(c.key[:], key64[:])
	return &c
}

func (c *Cipher) BlockSize() int { return 32 }

func (c *Cipher) Encrypt(dst, src []byte) {
	if len(src) != len(dst) || len(src) != 32 {
		panic("Cipher.Encrypt requires src and dst slices must both be exactly 32 bytes")
	}
	if &dst[0] != &src[0] {
		copy(dst, src)
	}
	encrypt256(convert.Inplace32BytesToUInt64(dst), &c.key, &c.tweak)
}

func (c *Cipher) Decrypt(dst, src []byte) {
	if len(src) != len(dst) || len(src) != 32 {
		panic("Cipher.Decrypt requires src and dst slices must both be exactly 32 bytes")
	}
	if &dst[0] != &src[0] {
		copy(dst, src)
	}
	decrypt256(convert.Inplace32BytesToUInt64(dst), &c.key, &c.tweak)
}

// TweakableBlockCipher implements types.TweakableBlockCipher.
type TweakableBlockCipher struct{}

func (t TweakableBlockCipher) Encrypt(data []byte, key []byte, tweak []byte) {
	Encrypt(data, key, tweak)
}

func (t TweakableBlockCipher) Decrypt(data []byte, key []byte, tweak []byte) {
	Decrypt(data, key, tweak)
}

func (t TweakableBlockCipher) BlockSize() int {
	return 256
}

func (t TweakableBlockCipher) TweakSize() int {
	return 128
}

// Encrypt encrypts a single block of 32 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 32 bytes of key, the rest is used internally.  Similarly,
// the tweak should be contained in the first two elements of tweak, the third value is used internally.
func Encrypt(data []byte, key []byte, tweak []byte) {
	if len(data) != 32 {
		panic("Encrypt requires that data is exactly 32 bytes")
	}
	if len(key) != 40 {
		panic("Encrypt requires that key is exactly 40 bytes")
	}
	encrypt256(convert.Inplace32BytesToUInt64(data), convert.Inplace40BytesToUInt64(key), convert.Inplace24BytesToUInt64(tweak))
}

// Decrypt decrypts a single block of 32 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 32 bytes of key, the rest is used internally.  Similarly,
// the tweak should be contained in the first two elements of tweak, the third value is used internally.
func Decrypt(data []byte, key []byte, tweak []byte) {
	if len(data) != 32 {
		panic("Decrypt requires that data is exactly 32 bytes")
	}
	if len(key) != 40 {
		panic("Encrypt requires that key is exactly 40 bytes")
	}
	decrypt256(convert.Inplace32BytesToUInt64(data), convert.Inplace40BytesToUInt64(key), convert.Inplace24BytesToUInt64(tweak))
}

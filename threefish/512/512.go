package threefish

import (
	"github.com/runningwild/skein/convert"
	"github.com/runningwild/skein/types"
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

// TweakableBlockCipher implements types.TweakableBlockCipher.
type TweakableBlockCipher struct{}

func (t TweakableBlockCipher) Encrypt(data []byte, key []byte, tweak []byte) {
	Encrypt(data, key, tweak)
}

func (t TweakableBlockCipher) Decrypt(data []byte, key []byte, tweak []byte) {
	Decrypt(data, key, tweak)
}

func (t TweakableBlockCipher) BlockSize() int {
	return 512
}

func (t TweakableBlockCipher) TweakSize() int {
	return 128
}

func (t TweakableBlockCipher) JFish() types.JFish {
	return MakeJFish()
}

// Encrypt encrypts a single block of 64 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 64 bytes of key, the rest is used internally.  Similarly,
// the tweak should be contained in the first two elements of tweak, the third value is used internally.
func Encrypt(data []byte, key []byte, tweak []byte) {
	if len(data) != 64 {
		panic("Encrypt requires that data is exactly 64 bytes")
	}
	if len(key) != 72 {
		panic("Encrypt requires that key is exactly 72 bytes")
	}
	if len(tweak) != 24 {
		panic("Encrypt requires that tweak is exactly 24 bytes")
	}
	encrypt512(convert.Inplace64BytesToUInt64(data), convert.Inplace72BytesToUInt64(key), convert.Inplace24BytesToUInt64(tweak))
}

// Decrypt decrypts a single block of 64 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 64 bytes of key, the rest is used internally.  Similarly,
// the tweak should be contained in the first two elements of tweak, the third value is used internally.
func Decrypt(data []byte, key []byte, tweak []byte) {
	if len(data) != 64 {
		panic("Decrypt requires that data is exactly 64 bytes")
	}
	if len(key) != 72 {
		panic("Decrypt requires that key is exactly 72 bytes")
	}
	if len(tweak) != 24 {
		panic("Decrypt requires that tweak is exactly 24 bytes")
	}
	decrypt512(convert.Inplace64BytesToUInt64(data), convert.Inplace72BytesToUInt64(key), convert.Inplace24BytesToUInt64(tweak))
}

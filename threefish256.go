package skein

type Cipher256 struct {
	key   [5]uint64 // The user-defined key is the first 4 values, we add the 5th.
	tweak [3]uint64 // The user-defined tweak is the first 2 values, we add the 3rd.
}

func MakeCipher256(key []byte) *Cipher256 {
	if len(key) != 32 {
		panic("MakeCipher256 can only be called with a key of exactly 32 bytes")
	}
	var b Cipher256
	var shortKey [4]uint64
	convert256InPlaceBytesToUint64(key, &shortKey)
	copy(b.key[:], shortKey[:])
	return &b
}

func (b *Cipher256) BlockSize() int { return 32 }

func (b *Cipher256) Encrypt(dst, src []byte) {
	if len(src) != len(dst) || len(src) != 32 {
		panic("Cipher256.Encrypt requires src and dst slices must both be exactly 32 bytes")
	}
	if &dst[0] != &src[0] {
		copy(dst, src)
	}
	Encrypt256(dst, &b.key, &b.tweak)
}

func (b *Cipher256) Decrypt(dst, src []byte) {
	if len(src) != len(dst) || len(src) != 32 {
		panic("Cipher256.Decrypt requires src and dst slices must both be exactly 32 bytes")
	}
	if &dst[0] != &src[0] {
		copy(dst, src)
	}
	decrypt256Simple(inplaceConvert32BytesToUInt64(dst), &b.key, &b.tweak)
}

// Encrypt256 encrypts a single block of 32 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 4 elements of key, the fifth value is used internally.
// Similarly, the tweak should be contained in the first two elements of tweak, the third value is
// used internally.
func Encrypt256(data []byte, key *[5]uint64, tweak *[3]uint64) {
	if len(data) != 32 {
		panic("Encrypt256 requires that data is exactly 32 bytes")
	}
	encrypt256(inplaceConvert32BytesToUInt64(data), key, tweak)
}

// Decrypt256 decrypts a single block of 32 bytes in place using key and tweak.  The actual key
// itself should be contained in the first 4 elements of key, the fifth value is used internally.
// Similarly, the tweak should be contained in the first two elements of tweak, the third value is
// used internally.
func Decrypt256(data []byte, key *[5]uint64, tweak *[3]uint64) {
	if len(data) != 32 {
		panic("Decrypt256 requires that data is exactly 32 bytes")
	}
	decrypt256(inplaceConvert32BytesToUInt64(data), key, tweak)
}

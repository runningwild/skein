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
		panic("Cipher256.Encrypt src and dst slices must both be exactly 32 bytes")
	}
	copy(dst, src)
	encrypt256Simple(inplaceCovertBytesToUInt64(dst), &b.key, &b.tweak)
}

func (b *Cipher256) Decrypt(dst, src []byte) {
	if len(src) != len(dst) || len(src) != 32 {
		panic("Cipher256.Encrypt src and dst slices must both be exactly 32 bytes")
	}
	copy(dst, src)
	decrypt256Simple(inplaceCovertBytesToUInt64(dst), &b.key, &b.tweak)
}

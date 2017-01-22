// +build amd64

package skein

type Block256 struct {
	state [4]uint64
	key   [5]uint64 // The user-defined key is the first 4 values, we add the 5th.
	tweak [3]uint64 // The user-defined tweak is the first 2 values, we add the 3rd.
}

func MakeBlock256(key [32]byte) *Block256 {
	var b Block256
	for i, v := range convert256BytesToUint64(key) {
		b.key[i] = v
	}
	return &b
}

func (b *Block256) BlockSize() int { return 32 }

func (b *Block256) Encrypt(dst, src []byte) {
	convert256InPlaceBytesToUint64(src, &b.state)
	encrypt256(&b.state, &b.key, &b.tweak)
	convert256InPlaceUint64ToBytes(&b.state, dst)
}

func (b *Block256) rawEncrypt(state *[4]uint64) {
	encrypt256(state, &b.key, &b.tweak)
}

func (b *Block256) Decrypt(dst, src []byte) {
	var buf [32]byte
	copy(buf[:], src)
	b.state = convert256BytesToUint64(buf)
	encrypt256(&b.state, &b.key, &b.tweak)
	buf = convert256Uint64ToBytes(b.state)
	copy(dst, buf[:])
}

func encrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64)
func decrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64)

package ubi

import (
	"fmt"
	"sync"

	"github.com/runningwild/skein/convert"
	"github.com/runningwild/skein/types"
)

func New(tbc types.TweakableBlockCipher) (*UBI, error) {
	if tbc.TweakSize() != 128 {
		return nil, fmt.Errorf("only tweak size 128 is supported")
	}

	switch tbc.BlockSize() {
	case 256:
	case 512:
	case 1024:
	default:
		return nil, fmt.Errorf("only block sizes 256, 512, and 1024 are supported")
	}

	return &UBI{
		tbc:          tbc,
		blockSize:    tbc.BlockSize(),
		blockBytes:   tbc.BlockSize() / 8,
		blockUint64s: tbc.BlockSize() / 64,
		gs:           make(map[uint64][]byte),
	}, nil
}

type UBI struct {
	tbc          types.TweakableBlockCipher
	blockSize    int
	blockBytes   int
	blockUint64s int

	mu sync.RWMutex
	gs map[uint64][]byte
}

func (ubi *UBI) UBI(G []byte, M []byte, Ts [2]uint64) []byte {
	return ubi.UBIBits(G, 0, M, Ts)
}

func (ubi *UBI) UBIBits(G []byte, lastByteBits int, M []byte, Ts [2]uint64) []byte {
	if lastByteBits < 0 || lastByteBits >= 8 {
		panic("lastByteBits must be in [0, 7]")
	}

	it := ubi.Iterate(G, Ts)
	for len(M) > ubi.blockBytes {
		it.Block(M[0:ubi.blockBytes])
		M = M[ubi.blockBytes:]
	}
	return it.Finish(M, lastByteBits)
}

type Iterator struct {
	ubi        *UBI
	tweak      []uint64
	tweakBytes []byte
	h          []byte
	buf        []byte
}

func (ubi *UBI) Iterate(G []byte, Ts [2]uint64) *Iterator {
	if len(G) != ubi.blockBytes {
		panic(fmt.Sprintf("G must match the block size, %d != %d", len(G), ubi.blockSize))
	}
	if Ts[1]&(1<<(119-64)) != 0 {
		panic("cannot call UBI with the BitPad field set on Ts")
	}
	if Ts[1]&(1<<(126-64)) != 0 {
		panic("cannot call UBI with the First field set on Ts")
	}
	if Ts[1]&(1<<(127-64)) != 0 {
		panic("cannot call UBI with the Final field set on Ts")
	}

	tweak := make([]uint64, 3)
	tweak[1] = Ts[1] | (1 << (126 - 64)) // set the 'first' bit

	// TODO: This is another case that isn't portable on big-endian machines.  Probably want to have
	// a tweak function in an unsafe package that can do this quickly on LE, and correctly on BE.
	tweakBytes := convert.Inplace3Uint64ToBytes(tweak)[:]

	H := make([]byte, len(G)+8)
	copy(H, G)

	buf := make([]byte, ubi.blockBytes)

	return &Iterator{
		ubi:        ubi,
		tweak:      tweak,
		tweakBytes: tweakBytes,
		h:          H,
		buf:        buf,
	}
}

func (it *Iterator) Block(M []byte) {
	copy(it.buf, M)

	// Here we aren't supporting sizes over 2^64, even though the spec supports up to 2^96.
	it.tweak[0] += uint64(it.ubi.blockBytes)

	it.ubi.tbc.Encrypt(it.buf, it.h, it.tweakBytes)
	convert.XorBytes(it.h[0:it.ubi.blockBytes], M, it.buf)
	it.tweak[1] &^= (1 << (126 - 64)) // unset the 'first' bit
}

func (it *Iterator) Finish(M []byte, lastByteBits int) []byte {
	var tweak [3]uint64
	copy(tweak[:], it.tweak)
	tweak[0] += uint64(len(M))
	tweak[1] |= (1 << (127 - 64)) // set the 'last' bit
	lastBlock := make([]byte, it.ubi.blockBytes)
	copy(lastBlock, M)
	if lastByteBits != 0 {
		tweak[1] |= (1 << (119 - 64)) // set the 'bitpad' bit
		b := lastBlock[len(M)-1]
		var lastUsedBit byte = 1 << uint(7-lastByteBits+1)
		b = (b &^ (lastUsedBit - 1)) | (lastUsedBit >> 1)
		lastBlock[len(M)-1] = b
	}
	buf := make([]byte, it.ubi.blockBytes)
	copy(buf, lastBlock)
	it.ubi.tbc.Encrypt(lastBlock, it.h, convert.Inplace3Uint64ToBytes(tweak[:])[:])
	convert.Xor(lastBlock, lastBlock, buf)
	return lastBlock
}

func (ubi *UBI) NewHasher(N int) *Hasher {
	h := &Hasher{
		ubi: ubi,
		buf: make([]byte, ubi.blockBytes)[0:0],
		n:   uint64(N),
	}
	h.Reset()
	return h
}

func (ubi *UBI) NewMACer(key []byte, N int) *Hasher {
	h := &Hasher{
		ubi: ubi,
		buf: make([]byte, ubi.blockBytes)[0:0],
		key: make([]byte, len(key)),
		n:   uint64(N),
	}
	copy(h.key, key)
	h.Reset()
	return h
}

func (ubi *UBI) Hash(M []byte, lastByteBits int, N uint64) []byte {
	return ubi.skein(nil, []Tuple{{TypeMsg, M, lastByteBits}}, N)
}

func (ubi *UBI) MAC(K []byte, M []byte, lastByteBits int, N uint64) []byte {
	return ubi.skein(K, []Tuple{{TypeMsg, M, lastByteBits}}, N)
}

type Tuple struct {
	Type         ConfigType
	Msg          []byte
	LastByteBits int
}

// Nb - The internal state size, this is known implicitly in the ubi object.
// No (N) - The output size, in bits.
// K - A key of Nk bytes. Set to the empty string (Nk = 0) if no key is desired.
// L List of t tuples (Ti,Mi) where Ti is a type value and Mi is a string of bits encoded in a string of bytes.
// NEXT: This method should be exported, and we shouldn't bother exporting the Hash and MAC methods above.  Or should we?
func (ubi *UBI) skein(K []byte, L []Tuple, N uint64) []byte {
	var Gn []byte = ubi.GetInitialChainingValue(K, N)
	for i := range L {
		if L[i].LastByteBits < 0 || L[i].LastByteBits > 7 {
			panic("LastByteBits must be in range [0, 7]")
		}
		Gn = ubi.UBIBits(Gn, L[i].LastByteBits, L[i].Msg, [2]uint64{0, uint64(L[i].Type)})
	}
	buf := make([]byte, int(N)/8+ubi.blockBytes)
	view := buf[:]
	// put c in an array so we can convert it to bytes to pass to UBI.
	var c [1]uint64
	cb := convert.Inplace1Uint64ToBytes(c[:])[:]
	iterations := (N + uint64(ubi.blockSize) - 1) / uint64(ubi.blockSize)
	for c[0] < iterations {
		copy(view, ubi.UBI(Gn, cb, tweakTypeOut))
		view = view[ubi.blockBytes:]
		c[0]++
	}
	if uint64(len(buf)*8) > N {
		buf = buf[0 : int(N+7)/8]
		if N&0x07 != 0 {
			// This masks away the upper bits that we don't care about, in the event that we asked for a
			// number of bits that doesn't evenly divide a byte.
			buf[len(buf)-1] = buf[len(buf)-1] & ((1 << uint(N&0x07)) - 1)
		}
	}
	return buf
}

func (ubi *UBI) GetInitialChainingValue(K []byte, N uint64) []byte {
	var G0 []byte
	if len(K) > 0 {
		Kcompressed := ubi.UBI(make([]byte, ubi.blockBytes), K, tweakTypeKey)
		G0 = ubi.UBI(Kcompressed, []byte{
			0x53, 0x48, 0x41, 0x33, // SHA3
			0x01, 0x00, // Version Number
			0x00, 0x00, // Reserved
			byte(N), byte(N >> 8), byte(N >> 16), byte(N >> 24), byte(N >> 32), byte(N >> 40), byte(N >> 48), byte(N >> 56), // Output size in bits (256)
			0x00, 0x00, 0x00, // Tree params
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Reserved,
		}, tweakTypeCfg)
	} else {
		var ok bool
		ubi.mu.RLock()
		G0, ok = ubi.gs[N]
		ubi.mu.RUnlock()
		if !ok {
			G0 = ubi.UBI(make([]byte, ubi.blockBytes), []byte{
				0x53, 0x48, 0x41, 0x33, // SHA3
				0x01, 0x00, // Version Number
				0x00, 0x00, // Reserved
				byte(N), byte(N >> 8), byte(N >> 16), byte(N >> 24), byte(N >> 32), byte(N >> 40), byte(N >> 48), byte(N >> 56), // Output size in bits (256)
				0x00, 0x00, 0x00, // Tree params
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Reserved,
			}, tweakTypeCfg)
			ubi.mu.Lock()
			ubi.gs[N] = G0
			ubi.mu.Unlock()
		}
	}
	return G0
}

type ConfigType uint64

const (
	TypeKey ConfigType = 0 << (120 - 64)
	TypeCfg ConfigType = 4 << (120 - 64)
	TypePrs ConfigType = 8 << (120 - 64)
	TypePK  ConfigType = 12 << (120 - 64)
	TypeKdf ConfigType = 16 << (120 - 64)
	TypeNon ConfigType = 20 << (120 - 64)
	TypeMsg ConfigType = 48 << (120 - 64)
	TypeOut ConfigType = 63 << (120 - 64)
)

var (
	tweakTypeKey = [2]uint64{0, uint64(TypeKey)}
	tweakTypeCfg = [2]uint64{0, uint64(TypeCfg)}
	tweakTypePrs = [2]uint64{0, uint64(TypePrs)}
	tweakTypePK  = [2]uint64{0, uint64(TypePK)}
	tweakTypeKdf = [2]uint64{0, uint64(TypeKdf)}
	tweakTypeNon = [2]uint64{0, uint64(TypeNon)}
	tweakTypeMsg = [2]uint64{0, uint64(TypeMsg)}
	tweakTypeOut = [2]uint64{0, uint64(TypeOut)}
)

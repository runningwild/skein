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

	state := ubi.start(G, Ts)
	for len(M) > ubi.blockBytes {
		ubi.block(state, M[0:ubi.blockBytes])
		M = M[ubi.blockBytes:]
	}
	return ubi.finish(state, M, lastByteBits)
}

type ubiState struct {
	tweak      []uint64
	tweakBytes []byte
	H          []byte
	buf        []byte
}

func (ubi *UBI) start(G []byte, Ts [2]uint64) *ubiState {
	tweak := make([]uint64, 3)
	tweak[1] = Ts[1] | (1 << (126 - 64)) // set the 'first' bit
	tweakBytes := convert.Inplace3Uint64ToBytes(tweak)[:]

	H := make([]byte, len(G)+8)
	copy(H, G)

	buf := make([]byte, ubi.blockBytes)

	return &ubiState{
		tweak:      tweak,
		tweakBytes: tweakBytes,
		H:          H,
		buf:        buf,
	}
}

func (ubi *UBI) block(state *ubiState, M []byte) {
	copy(state.buf, M)

	// Here we aren't supporting sizes over 2^64, even though the spec supports up to 2^96.
	state.tweak[0] += uint64(ubi.blockBytes)

	ubi.tbc.Encrypt(state.buf, state.H, state.tweakBytes)
	convert.XorBytes(state.H[0:ubi.blockBytes], M, state.buf)
	state.tweak[1] &^= (1 << (126 - 64)) // unset the 'first' bit
}

func (ubi *UBI) finish(state *ubiState, M []byte, lastByteBits int) []byte {
	state.tweak[0] += uint64(len(M))
	state.tweak[1] |= (1 << (127 - 64)) // set the 'last' bit
	lastBlock := make([]byte, ubi.blockBytes)
	copy(lastBlock, M)
	if lastByteBits != 0 {
		state.tweak[1] |= (1 << (119 - 64)) // set the 'bitpad' bit
		b := lastBlock[len(M)-1]
		var lastUsedBit byte = 1 << uint(7-lastByteBits+1)
		b = (b &^ (lastUsedBit - 1)) | (lastUsedBit >> 1)
		lastBlock[len(M)-1] = b
	}
	copy(state.buf, lastBlock)
	ubi.tbc.Encrypt(lastBlock, state.H, state.tweakBytes)
	convert.Xor(state.H[0:ubi.blockBytes], lastBlock, state.buf)
	return state.H[0:ubi.blockBytes]
}

func (ubi *UBI) Hash(M []byte, MBits int, N uint64) []byte {
	return ubi.skein(nil, []tuple{{typeMsg, M, MBits}}, N)
}

func (ubi *UBI) MAC(K []byte, M []byte, MBits int, N uint64) []byte {
	return ubi.skein(K, []tuple{{typeMsg, M, MBits}}, N)
}

type tuple struct {
	typ     configType
	msg     []byte
	msgBits int
}

// Nb - The internal state size, this is known implicitly in the ubi object.
// No (N) - The output size, in bits.
// K - A key of Nk bytes. Set to the empty string (Nk = 0) if no key is desired.
// L List of t tuples (Ti,Mi) where Ti is a type value and Mi is a string of bits encoded in a string of bytes.
func (ubi *UBI) skein(K []byte, L []tuple, N uint64) []byte {
	for i := range L {
		if len(L[i].msg) != (L[i].msgBits+7)/8 {
			panic(fmt.Sprintf("len(L[%d].msg) and L[%d].msgBits do not match", i, i))
		}
	}
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
		}, [2]uint64{0, 4 << (120 - 64)})
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
			}, [2]uint64{0, 4 << (120 - 64)})
			ubi.mu.Lock()
			ubi.gs[N] = G0
			ubi.mu.Unlock()
		}
	}
	var Gn []byte = G0
	for i := range L {
		Gn = ubi.UBIBits(Gn, L[i].msgBits&0x07, L[i].msg, [2]uint64{0, uint64(L[i].typ)})
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

var (
	tweakTypeKey = [2]uint64{0, uint64(typeKey)}
	tweakTypeMsg = [2]uint64{0, uint64(typeMsg)}
	tweakTypeOut = [2]uint64{0, uint64(typeOut)}
)

func ConfigTweak(position uint64, treeLevel uint64, bitPad bool, typ configType, first, final bool) [2]uint64 {
	var block [2]uint64
	block[0] = position
	if bitPad {
		block[1] |= (1 << (119 - 64))
	}
	block[1] = uint64(typ)
	if first {
		block[1] |= (1 << (126 - 64))
	}
	if final {
		block[1] |= (1 << (127 - 64))
	}
	return block
}

type configType uint64

const (
	typeKey configType = 0 << (120 - 64)
	typeCfg configType = 4 << (120 - 64)
	typePrs configType = 8 << (120 - 64)
	typePK  configType = 12 << (120 - 64)
	typeKdf configType = 16 << (120 - 64)
	typeNon configType = 20 << (120 - 64)
	typeMsg configType = 48 << (120 - 64)
	typeOut configType = 63 << (120 - 64)
)

package ubi

import (
	"fmt"
	"hash"
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
	return it.FinishSafe(M, lastByteBits)
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
	it.tweak[0] += uint64(len(M))
	it.tweak[1] |= (1 << (127 - 64)) // set the 'last' bit
	lastBlock := make([]byte, it.ubi.blockBytes)
	copy(lastBlock, M)
	if lastByteBits != 0 {
		it.tweak[1] |= (1 << (119 - 64)) // set the 'bitpad' bit
		b := lastBlock[len(M)-1]
		var lastUsedBit byte = 1 << uint(7-lastByteBits+1)
		b = (b &^ (lastUsedBit - 1)) | (lastUsedBit >> 1)
		lastBlock[len(M)-1] = b
	}
	copy(it.buf, lastBlock)
	it.ubi.tbc.Encrypt(lastBlock, it.h, it.tweakBytes)
	convert.Xor(it.h[0:it.ubi.blockBytes], lastBlock, it.buf)
	return it.h[0:it.ubi.blockBytes]
}

func (it *Iterator) FinishSafe(M []byte, lastByteBits int) []byte {
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

type hasher struct {
	ubi *UBI
	buf []byte
	key []byte
	n   uint64
	it  *Iterator
}

func (ubi *UBI) NewHasher(N int) hash.Hash {
	h := &hasher{
		ubi: ubi,
		buf: make([]byte, ubi.blockBytes)[0:0],
		n:   uint64(N),
	}
	h.Reset()
	return h
}

func (h *hasher) Write(b []byte) (n int, err error) {
	written := len(b)

	// If the buffer plus b won't fill up more than a full block, just add it to the buffer and return.
	if len(h.buf)+len(b) <= h.ubi.blockBytes {
		h.buf = append(h.buf, b...)
		return written, nil
	}

	// If we have something in the buffer we need to make a full block using that first.
	if len(h.buf) > 0 {
		amt := h.ubi.blockBytes - len(h.buf)
		h.it.Block(append(h.buf, b[0:amt]...))
		h.buf = h.buf[0:0]
		b = b[amt:]
	}

	// Iterate until we don't have more than a full block.
	for len(b) > h.ubi.blockBytes {
		h.it.Block(b[0:h.ubi.blockBytes])
		b = b[h.ubi.blockBytes:]
	}

	// Store anything up to a single remaining full block for later.
	h.buf = append(h.buf, b...)

	return written, nil
}
func (h *hasher) Sum(b []byte) []byte {
	Gn := h.it.FinishSafe(h.buf, 0)

	buf := make([]byte, int(h.n)/8+h.ubi.blockBytes)
	view := buf[:]
	// put c in an array so we can convert it to bytes to pass to UBI.
	var c [1]uint64
	cb := convert.Inplace1Uint64ToBytes(c[:])[:]
	iterations := (h.n + uint64(h.ubi.blockSize) - 1) / uint64(h.ubi.blockSize)
	for c[0] < iterations {
		copy(view, h.ubi.UBI(Gn, cb, tweakTypeOut))
		view = view[h.ubi.blockBytes:]
		c[0]++
	}
	if uint64(len(buf)*8) > h.n {
		buf = buf[0 : int(h.n+7)/8]
		if h.n&0x07 != 0 {
			// This masks away the upper bits that we don't care about, in the event that we asked for a
			// number of bits that doesn't evenly divide a byte.
			buf[len(buf)-1] = buf[len(buf)-1] & ((1 << uint(h.n&0x07)) - 1)
		}
	}
	return append(b, buf...)
}
func (h *hasher) Reset() {
	Gn := h.ubi.getInitialChainingValue(h.key, h.n)
	h.it = h.ubi.Iterate(Gn, tweakTypeMsg)
	h.buf = h.buf[0:0]
}
func (h *hasher) Size() int {
	return (int(h.n) + 7) / 8
}
func (h *hasher) BlockSize() int {
	return h.ubi.blockSize
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
	var Gn []byte = ubi.getInitialChainingValue(K, N)
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

func (ubi *UBI) getInitialChainingValue(K []byte, N uint64) []byte {
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

var (
	tweakTypeKey = [2]uint64{0, uint64(typeKey)}
	tweakTypeCfg = [2]uint64{0, uint64(typeCfg)}
	tweakTypePrs = [2]uint64{0, uint64(typePrs)}
	tweakTypePK  = [2]uint64{0, uint64(typePK)}
	tweakTypeKdf = [2]uint64{0, uint64(typeKdf)}
	tweakTypeNon = [2]uint64{0, uint64(typeNon)}
	tweakTypeMsg = [2]uint64{0, uint64(typeMsg)}
	tweakTypeOut = [2]uint64{0, uint64(typeOut)}
)

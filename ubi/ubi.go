package ubi

import (
	"fmt"

	"github.com/runningwild/skein/convert"
)

type TweakableBlockCipher func(data []byte, key []byte, tweak *[3]uint64)

func New(tbc TweakableBlockCipher, blockSize int) (*UBI, error) {
	if blockSize <= 0 || blockSize%8 != 0 {
		return nil, fmt.Errorf("blockSize must be a positive multiple of 8")
	}

	var convertBlockBytesToUint64 func([]byte) []uint64
	var convertBlockUint64ToBytes func([]uint64) []byte
	switch blockSize {
	case 256:
		convertBlockBytesToUint64 = func(b []byte) []uint64 {
			return convert.Inplace32BytesToUInt64(b)[:]
		}
		convertBlockUint64ToBytes = func(v []uint64) []byte {
			return convert.Inplace4Uint64ToBytes(v)[:]
		}

	case 512:
		convertBlockBytesToUint64 = func(b []byte) []uint64 {
			return convert.Inplace64BytesToUInt64(b)[:]
		}
		convertBlockUint64ToBytes = func(v []uint64) []byte {
			return convert.Inplace8Uint64ToBytes(v)[:]
		}

	case 1024:
		convertBlockBytesToUint64 = func(b []byte) []uint64 {
			return convert.Inplace128BytesToUInt64(b)[:]
		}
		convertBlockUint64ToBytes = func(v []uint64) []byte {
			return convert.Inplace16Uint64ToBytes(v)[:]
		}

	default:
		return nil, fmt.Errorf("only block sizes 256, 512, and 1024 are supported")
	}

	return &UBI{
		tbc:                       tbc,
		blockSize:                 blockSize,
		blockBytes:                blockSize / 8,
		blockUint64s:              blockSize / 64,
		convertBlockBytesToUint64: convertBlockBytesToUint64,
		convertBlockUint64ToBytes: convertBlockUint64ToBytes,
		Gs: make(map[uint64][]byte),
	}, nil
}

type UBI struct {
	tbc          TweakableBlockCipher
	blockSize    int
	blockBytes   int
	blockUint64s int

	convertBlockBytesToUint64 func([]byte) []uint64
	convertBlockUint64ToBytes func([]uint64) []byte

	Gs map[uint64][]byte
}

func (ubi *UBI) UBI(G []byte, M []byte, Ts [2]uint64) []byte {
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

	var tweak [3]uint64
	tweak[1] = Ts[1] | (1 << (126 - 64)) // set the 'first' bit

	H := make([]byte, len(G)+8)
	copy(H, G)

	// Figure out how much belongs in the 'last' block.  There must be a last block, even if it's
	// zero bits.
	rem := len(M) % ubi.blockBytes
	if rem == 0 {
		rem = ubi.blockBytes
		if rem > len(M) {
			rem = len(M)
		}
	}
	lastBlock := make([]byte, ubi.blockBytes)
	copy(lastBlock, M[len(M)-rem:])
	M = M[0 : len(M)-rem]

	state64 := make([]uint64, ubi.blockUint64s)
	state := ubi.convertBlockUint64ToBytes(state64)
	// Process every full block except the last.
	for len(M) > 0 {
		M64 := ubi.convertBlockBytesToUint64(M[0:ubi.blockBytes])
		copy(state64, M64)

		// Here we aren't supporting sizes over 2^64, even though the spec supports up to 2^96.
		tweak[0] += uint64(ubi.blockBytes)

		ubi.tbc(state, H, &tweak)
		convert.Xor(H[0:ubi.blockBytes], M[0:ubi.blockBytes], state)
		M = M[ubi.blockBytes:]
		tweak[1] &^= (1 << (126 - 64)) // unset the 'first' bit
	}

	// Process the last block.
	tweak[0] += uint64(rem)
	tweak[1] |= (1 << (127 - 64)) // set the 'last' bit
	block64 := ubi.convertBlockBytesToUint64(lastBlock)
	copy(state64, block64)
	ubi.tbc(lastBlock, H, &tweak)
	convert.Xor(H[0:ubi.blockBytes], lastBlock, state)
	return H[0:ubi.blockBytes]
}

func (ubi *UBI) Skein(M []byte, msgLen int, N uint64) []byte {
	G0, ok := ubi.Gs[N]
	if !ok {
		G0 = ubi.UBI(make([]byte, ubi.blockBytes), []byte{
			0x53, 0x48, 0x41, 0x33, // SHA3
			0x01, 0x00, // Version Number
			0x00, 0x00, // Reserved
			byte(N), byte(N >> 8), byte(N >> 16), byte(N >> 24), byte(N >> 32), byte(N >> 40), byte(N >> 48), byte(N >> 56), // Output size in bits (256)
			0x00, 0x00, 0x00, // Tree params
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Reserved,
		}, [2]uint64{0, 4 << (120 - 64)})
		ubi.Gs[N] = G0
	}
	G1 := ubi.UBI(G0, M, tweakTypeMsg)
	buf := make([]byte, int(N)/8+ubi.blockBytes)
	view := buf[:]
	// put c in an array so we can convert it to bytes to pass to UBI.
	var c [1]uint64
	cb := convert.Inplace1Uint64ToBytes(c[:])[:]
	iterations := (N + uint64(ubi.blockSize) - 1) / uint64(ubi.blockSize)
	for c[0] < iterations {
		copy(view, ubi.UBI(G1, cb, tweakTypeOut))
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

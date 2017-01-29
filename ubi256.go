package skein

import "fmt"

func ubi256(G [4]uint64, M []byte, Ts [2]uint64) [4]uint64 {
	pos := 0

	var key [5]uint64
	var tweak [3]uint64

	var first uint64 = 1 << (126 - 64)
	var last uint64 = 0
	H := G

	// Annoying thing in the loop condition is to make sure we run this once for a zero-bit message.
	start, pos := 0, 0
	var buf [32]byte
	for last == 0 {
		if start+32 >= len(M) {
			last = 1 << (127 - 64)
			pos = len(M)
			for i := start; i < pos; i++ {
				buf[i-start] = M[i]
			}
			for i := pos - start; i < 32; i++ {
				buf[i] = 0
			}
		} else {
			pos += 32
			copy(buf[:], M[start:pos])
		}
		tweak[1] = Ts[1] | first | last
		first = 0

		// Here we aren't supporting sizes over 2^64, even though the spec supports up to 2^96.
		tweak[0] = Ts[0] + uint64(pos)

		start = pos

		copy(key[:], H[:])
		msg64 := inplaceCovertBytesToUInt64(buf[:])
		state := *msg64
		encrypt256(&state, &key, &tweak)
		for i := range H {
			H[i] = msg64[i] ^ state[i]
		}
	}
	return H
}

type Skein256PRGN struct {
	state [4]uint64
}

func (s *Skein256PRGN) Seed(state [4]uint64) {
	s.state = state
}

func (s *Skein256PRGN) Read(b []byte) (n int, err error) {
	nextState := ubi256(s.state,
		[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		tweakTypeOut)
	var c uint64
	for c*32 < uint64(len(b)) {
		H := ubi256(s.state,
			[]byte{byte(c), byte(c >> 8), byte(c >> 16), byte(c >> 24), byte(c >> 32), byte(c >> 40), byte(c >> 48), byte(c >> 56)},
			tweakTypeOut)
		Hbytes := convert256Uint64ToBytes(H)
		copy(b[int(c)*32:], Hbytes[:])
		c++
	}
	s.state = nextState
	return len(b), nil
}

func Skein256_128(M []byte) [16]byte {
	G0 := ubi256([4]uint64{}, []byte{
		0x53, 0x48, 0x41, 0x33, // SHA1
		0x01, 0x00, // Version Number
		0x00, 0x00, // Reserved
		0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Output size in bits (128)
		0x00, 0x00, 0x00, // Tree params
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Reserved,
	}, [2]uint64{0, 4 << (120 - 64)})
	panic(fmt.Sprintf("%X", G0))
	G1 := ubi256(G0, M, [2]uint64{0, 48 << (120 - 64)})
	H := ubi256(G1, []byte{0, 0, 0, 0, 0, 0, 0, 0}, [2]uint64{0, 63 << (120 - 64)})
	var buf [16]byte
	Hbytes := convert256Uint64ToBytes(H)
	copy(buf[:], Hbytes[:])
	return buf
}

var (
	skein256_256_cfg = [4]uint64{0xFC9DA860D048B449, 0x2FCA66479FA7D833, 0xB33BC3896656840F, 0x6A54E920FDE8DA69}

	emptyMsg     = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	tweakTypeMsg = [2]uint64{0, typeMsg}
	tweakTypeOut = [2]uint64{0, typeOut}
)

const (
	typeMsg = 48 << (120 - 64)
	typeOut = 63 << (120 - 64)
)

func Skein256_256(M []byte) [32]byte {
	G1 := ubi256(skein256_256_cfg, M, tweakTypeMsg)
	H := ubi256(G1, emptyMsg, tweakTypeOut)
	return convert256Uint64ToBytes(H)
}

func Skein256_N(M []byte, N uint64) []byte {
	G0 := ubi256([4]uint64{}, []byte{
		0x53, 0x48, 0x41, 0x33, // SHA1
		0x01, 0x00, // Version Number
		0x00, 0x00, // Reserved
		byte(N), byte(N >> 8), byte(N >> 16), byte(N >> 24), byte(N >> 32), byte(N >> 40), byte(N >> 48), byte(N >> 56), // Output size in bits (256)
		0x00, 0x00, 0x00, // Tree params
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Reserved,
	}, [2]uint64{0, 4 << (120 - 64)})
	G1 := ubi256(G0, M, [2]uint64{0, 48 << (120 - 64)})
	buf := make([]byte, int(N)+32)
	var c uint64
	for c*32 < N {
		H := ubi256(G1,
			[]byte{byte(c), byte(c >> 8), byte(c >> 16), byte(c >> 24), byte(c >> 32), byte(c >> 40), byte(c >> 48), byte(c >> 56)},
			[2]uint64{0, 63 << (120 - 64)})
		Hbytes := convert256Uint64ToBytes(H)
		copy(buf[int(c)*32:int(c+1)*32], Hbytes[:])
		c++
	}
	if uint64(len(buf)*8) > N {
		buf = buf[0 : int(N+7)/8]
	}
	if N%8 != 0 {
		// This masks away the upper bits that we don't care about, in the event that we asked for a
		// number of bits that doesn't evenly divide a byte.
		buf[len(buf)-1] = buf[len(buf)-1] & ((1 << uint(N%8)) - 1)
	}
	return buf
}

func convert256Uint64ToBytes(v [4]uint64) [32]byte {
	var b [32]byte
	for i := range v {
		x := i * 8
		b[x] = byte(v[i])
		b[x+1] = byte(v[i] >> 8)
		b[x+2] = byte(v[i] >> 16)
		b[x+3] = byte(v[i] >> 24)
		b[x+4] = byte(v[i] >> 32)
		b[x+5] = byte(v[i] >> 40)
		b[x+6] = byte(v[i] >> 48)
		b[x+7] = byte(v[i] >> 56)
	}
	return b
}

func convert256InPlaceUint64ToBytes(v *[4]uint64, b []byte) {
	for i := range v {
		x := i * 8
		b[x] = byte(v[i])
		b[x+1] = byte(v[i] >> 8)
		b[x+2] = byte(v[i] >> 16)
		b[x+3] = byte(v[i] >> 24)
		b[x+4] = byte(v[i] >> 32)
		b[x+5] = byte(v[i] >> 40)
		b[x+6] = byte(v[i] >> 48)
		b[x+7] = byte(v[i] >> 56)
	}
}
func convert256InPlaceBytesToUint64(b []byte, v *[4]uint64) {
	for i := range v {
		x := i * 8
		v[i] = uint64(b[x]) |
			(uint64(b[x+1]) << 8) |
			(uint64(b[x+2]) << 16) |
			(uint64(b[x+3]) << 24) |
			(uint64(b[x+4]) << 32) |
			(uint64(b[x+5]) << 40) |
			(uint64(b[x+6]) << 48) |
			(uint64(b[x+7]) << 56)
	}
}

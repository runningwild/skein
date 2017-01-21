package skein

func ubi256(G [4]uint64, M []byte, Ts [2]uint64) [4]uint64 {
	pos := 0

	var b block256
	H := G
	// Annoying thing in the loop condition is to make sure we run this once for a zero-bit message.
	for i := 0; i < len(M) || (i == 0 && len(M) == 0); i += 32 {
		b.tweak[1] = Ts[1]
		if i == 0 {
			// Set the 'first' bit.
			b.tweak[1] = b.tweak[1] | (1 << (126 - 64))
		}
		if i+32 < len(M) {
			pos += 32
		} else {
			pos += len(M) - i
			// Set the 'final' bit.
			b.tweak[1] = b.tweak[1] | (1 << (127 - 64))
		}
		// Here we aren't supporting sizes over 2^64, even though the spec supports up to 2^96.
		b.tweak[0] = Ts[0] + uint64(pos)

		var msg64 [4]uint64
		for j := 0; j+i < len(M) && j < 32; j++ {
			msg64[j/8] += uint64(M[j+i]) << uint((j%8)*8)
		}
		copy(b.key[:], H[:])
		b.state = msg64
		b.encrypt()
		for i := range H {
			H[i] = msg64[i] ^ b.state[i]
		}
	}
	return H
}

func skein256_128(M []byte) [16]byte {
	G0 := ubi256([4]uint64{}, []byte{
		0x53, 0x48, 0x41, 0x33, // SHA1
		0x01, 0x00, // Version Number
		0x00, 0x00, // Reserved
		0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Output size in bits (128)
		0x00, 0x00, 0x00, // Tree params
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Reserved,
	}, [2]uint64{0, 4 << (120 - 64)})
	G1 := ubi256(G0, M, [2]uint64{0, 48 << (120 - 64)})
	H := ubi256(G1, []byte{0, 0, 0, 0, 0, 0, 0, 0}, [2]uint64{0, 63 << (120 - 64)})
	var buf [16]byte
	Hbytes := convert256Uint64ToBytes(H)
	copy(buf[:], Hbytes[:])
	return buf
}

func skein256_256(M []byte) [32]byte {
	G0 := ubi256([4]uint64{}, []byte{
		0x53, 0x48, 0x41, 0x33, // SHA1
		0x01, 0x00, // Version Number
		0x00, 0x00, // Reserved
		0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Output size in bits (256)
		0x00, 0x00, 0x00, // Tree params
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Reserved,
	}, [2]uint64{0, 4 << (120 - 64)})
	G1 := ubi256(G0, M, [2]uint64{0, 48 << (120 - 64)})
	H := ubi256(G1, []byte{0, 0, 0, 0, 0, 0, 0, 0}, [2]uint64{0, 63 << (120 - 64)})
	Hbytes := convert256Uint64ToBytes(H)
	var buf [32]byte
	copy(buf[:], Hbytes[:])
	return buf
}

func skein256_N(M []byte, N uint64) []byte {
	G0 := ubi256([4]uint64{}, []byte{
		0x53, 0x48, 0x41, 0x33, // SHA1
		0x01, 0x00, // Version Number
		0x00, 0x00, // Reserved
		byte(N), byte(N >> 8), byte(N >> 16), byte(N >> 24), byte(N >> 32), byte(N >> 40), byte(N >> 48), byte(N >> 56), // Output size in bits (256)
		0x00, 0x00, 0x00, // Tree params
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Reserved,
	}, [2]uint64{0, 4 << (120 - 64)})
	G1 := ubi256(G0, M, [2]uint64{0, 48 << (120 - 64)})
	var buf []byte
	var c uint64
	for uint64(len(buf)*8) < N {
		H := ubi256(G1,
			[]byte{byte(c), byte(c >> 8), byte(c >> 16), byte(c >> 24), byte(c >> 32), byte(c >> 40), byte(c >> 48), byte(c >> 56)},
			[2]uint64{0, 63 << (120 - 64)})
		Hbytes := convert256Uint64ToBytes(H)
		buf = append(buf, Hbytes[:]...)
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
	for i := range b {
		b[i] = byte(v[i/8] >> (uint(i%8) * 8))
	}
	return b
}

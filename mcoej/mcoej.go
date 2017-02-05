package mcoej

import (
	"fmt"

	"github.com/runningwild/skein/convert"
	"github.com/runningwild/skein/ubi"
)

func New(enc, dec TweakableBlockCipher, blockSize int) (*McOEJ, error) {
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

	skein, err := ubi.New(ubi.TweakableBlockCipher(enc), blockSize)
	if err != nil {
		return nil, fmt.Errorf("unable to make UBI: %v", err)
	}

	return &McOEJ{
		enc:          enc,
		dec:          dec,
		blockSize:    blockSize,
		blockBytes:   blockSize / 8,
		blockUint64s: blockSize / 64,
		skein:        skein,
		convertBlockBytesToUint64: convertBlockBytesToUint64,
		convertBlockUint64ToBytes: convertBlockUint64ToBytes,
	}, nil
}

type TweakableBlockCipher func(data []byte, key []byte, tweak *[3]uint64)

type McOEJ struct {
	enc, dec     TweakableBlockCipher
	blockSize    int
	blockBytes   int
	blockUint64s int
	skein        *ubi.UBI

	convertBlockBytesToUint64 func([]byte) []uint64
	convertBlockUint64ToBytes func([]uint64) []byte
}

func (j *McOEJ) Lock(key, nonce, publicData, plaintext, dst []byte) []byte {
	if len(key) != j.blockBytes {
		panic("key length must equal block size")
	}
	if len(nonce) != j.blockBytes {
		panic("nonce length must equal block size")
	}
	mac := j.skein.MAC(key, publicData, len(publicData)*8, uint64(j.blockSize))
	for i := range mac {
		mac[i] ^= nonce[i]
	}
	var fullTweak []uint64
	for len(mac) > 0 {
		fullTweak = append(fullTweak, packBytesIntoUint64(mac[0:8]))
		mac = mac[8:]
	}
	fullKey := make([]byte, len(key)+8)
	copy(fullKey, key)
	var tweak [3]uint64
	target := make([]byte, j.blockBytes)
	lastBlockLen := 0
	for len(plaintext) > 0 {
		lastBlockLen = j.blockBytes
		if len(plaintext) < len(target) {
			lastBlockLen = len(plaintext)
			for i := range target {
				target[i] = 0
			}
		}
		copy(target, plaintext)
		copy(tweak[0:2], fullTweak)
		fullTweak = append(fullTweak[2:], fullTweak[0:2]...)
		for i := 0; i < len(fullTweak); i++ {
			fullTweak[i] ^= packBytesIntoUint64(target[i*8 : (i+1)*8])
		}
		j.enc(target, fullKey, &tweak)
		for i := 0; i < len(fullTweak); i++ {
			fullTweak[i] ^= packBytesIntoUint64(target[i*8 : (i+1)*8])
		}
		dst = append(dst, target...)
		plaintext = plaintext[lastBlockLen:]
	}
	copy(tweak[0:2], fullTweak)
	fullTweak = fullTweak[2:]
	lastBlock := make([]byte, j.blockBytes)
	for i, v := range fullTweak {
		copy(lastBlock[i*8:], unpackBytesFromUint64(v))
	}
	lastBlock[len(fullTweak)*8] = byte(lastBlockLen)
	j.enc(lastBlock, fullKey, &tweak)
	dst = append(dst, lastBlock...)
	return dst
}

func (j *McOEJ) Unlock(key, nonce, publicData, ciphertext, dst []byte) ([]byte, error) {
	if len(key) != j.blockBytes {
		panic("key length must equal block size")
	}
	if len(nonce) != j.blockBytes {
		panic("nonce length must equal block size")
	}
	mac := j.skein.MAC(key, publicData, len(publicData)*8, uint64(j.blockSize))
	for i := range mac {
		mac[i] ^= nonce[i]
	}
	var fullTweak []uint64
	for len(mac) > 0 {
		fullTweak = append(fullTweak, packBytesIntoUint64(mac[0:8]))
		mac = mac[8:]
	}
	fullKey := make([]byte, len(key)+8)
	copy(fullKey, key)
	var tweak [3]uint64
	target := make([]byte, j.blockBytes)
	for len(ciphertext) > 0 {
		copy(target, ciphertext)
		ciphertext = ciphertext[j.blockBytes:]
		copy(tweak[0:2], fullTweak)
		if len(ciphertext) == 0 {
			fullTweak = fullTweak[2:]
		} else {
			fullTweak = append(fullTweak[2:], fullTweak[0:2]...)
		}
		if len(ciphertext) > 0 {
			for i := 0; i < len(fullTweak); i++ {
				fullTweak[i] ^= packBytesIntoUint64(target[i*8 : (i+1)*8])
			}
		}
		j.dec(target, fullKey, &tweak)
		if len(ciphertext) > 0 {
			for i := 0; i < len(fullTweak); i++ {
				fullTweak[i] ^= packBytesIntoUint64(target[i*8 : (i+1)*8])
			}
		}
		dst = append(dst, target...)
	}
	lastBlock := dst[len(dst)-j.blockBytes:]
	bad := false
	for i := range fullTweak {
		if fullTweak[i] != packBytesIntoUint64(lastBlock[i*8:(i+1)*8]) {
			bad = true
		}
	}
	if bad {
		return nil, fmt.Errorf("authentication failed")
	}
	lastBlockLen := int(lastBlock[len(fullTweak)*8])
	if lastBlockLen > j.blockBytes {
		return nil, fmt.Errorf("authentication failed")
	}
	if lastBlockLen == 0 {
		lastBlockLen = j.blockBytes
	}
	dst = dst[0 : len(dst)-2*j.blockBytes+lastBlockLen]
	return dst, nil
}

func packBytesIntoUint64(b []byte) uint64 {
	return uint64(b[0]) |
		(uint64(b[1]) << 8) |
		(uint64(b[2]) << 16) |
		(uint64(b[3]) << 24) |
		(uint64(b[4]) << 32) |
		(uint64(b[5]) << 40) |
		(uint64(b[6]) << 48) |
		(uint64(b[7]) << 56)
}

func unpackBytesFromUint64(u uint64) []byte {
	return []byte{
		byte(u),
		byte(u >> 8),
		byte(u >> 16),
		byte(u >> 24),
		byte(u >> 32),
		byte(u >> 40),
		byte(u >> 48),
		byte(u >> 56),
	}
}

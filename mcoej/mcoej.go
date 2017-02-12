package mcoej

import (
	"fmt"

	"github.com/runningwild/skein/convert"
	"github.com/runningwild/skein/ubi"
)

func New(enc, dec TweakableBlockCipher, blockSize int) (*McOEJ, error) {
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
		enc:            enc,
		dec:            dec,
		blockSize:      blockSize,
		blockBytes:     blockSize / 8,
		blockBytesMask: blockSize/8 - 1,
		blockUint64s:   blockSize / 64,
		skein:          skein,
		convertBlockBytesToUint64: convertBlockBytesToUint64,
		convertBlockUint64ToBytes: convertBlockUint64ToBytes,
	}, nil
}

type TweakableBlockCipher func(data []byte, key []byte, tweak []byte)

type McOEJ struct {
	enc, dec       TweakableBlockCipher
	blockSize      int
	blockBytes     int
	blockBytesMask int
	blockUint64s   int
	skein          *ubi.UBI

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
	fullTweak := j.skein.Hash(publicData, len(publicData)*8, uint64(j.blockSize))
	convert.Xor(fullTweak, fullTweak, nonce)
	tweakPos := 0
	fullKey := make([]byte, len(key)+8)
	copy(fullKey, key)
	var tweak [24]byte
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
		copy(tweak[0:16], fullTweak[tweakPos:])
		tweakPos = (tweakPos + 16) & j.blockBytesMask
		convert.Xor(fullTweak, fullTweak, target)
		j.enc(target, fullKey, tweak[:])
		convert.Xor(fullTweak, fullTweak, target)
		dst = append(dst, target...)
		plaintext = plaintext[lastBlockLen:]
	}
	copy(tweak[0:16], fullTweak[tweakPos:])
	fullTweak = fullTweak[16:]
	lastBlock := make([]byte, j.blockBytes)
	copy(lastBlock, fullTweak)
	lastBlock[len(fullTweak)] = byte(lastBlockLen)
	j.enc(lastBlock, fullKey, tweak[:])
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
	fullTweak := j.skein.Hash(publicData, len(publicData)*8, uint64(j.blockSize))
	convert.Xor(fullTweak, fullTweak, nonce)
	tweakPos := 0
	fullKey := make([]byte, len(key)+8)
	copy(fullKey, key)
	var tweak [24]byte
	target := make([]byte, j.blockBytes)
	for len(ciphertext) > 0 {
		copy(target, ciphertext)
		ciphertext = ciphertext[j.blockBytes:]
		copy(tweak[0:16], fullTweak[tweakPos:])
		if len(ciphertext) == 0 {
			fullTweak = fullTweak[16:]
		}
		tweakPos = (tweakPos + 16) & j.blockBytesMask
		if len(ciphertext) > 0 {
			convert.Xor(fullTweak, fullTweak, target)
		}
		j.dec(target, fullKey, tweak[:])
		if len(ciphertext) > 0 {
			convert.Xor(fullTweak, fullTweak, target)
		}
		dst = append(dst, target...)
	}
	lastBlock := dst[len(dst)-j.blockBytes:]
	bad := false
	for i := range fullTweak {
		if fullTweak[i] != lastBlock[i] {
			bad = true
		}
	}
	if bad {
		return nil, fmt.Errorf("authentication failed")
	}
	lastBlockLen := int(lastBlock[len(fullTweak)])
	if lastBlockLen > j.blockBytes {
		return nil, fmt.Errorf("authentication failed")
	}
	if lastBlockLen == 0 {
		lastBlockLen = j.blockBytes
	}
	dst = dst[0 : len(dst)-2*j.blockBytes+lastBlockLen]
	return dst, nil
}

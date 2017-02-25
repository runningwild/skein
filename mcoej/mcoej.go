package mcoej

import (
	"crypto/subtle"
	"fmt"

	"github.com/runningwild/skein/convert"
	"github.com/runningwild/skein/types"
)

func New(tbc types.TweakableBlockCipher) (*McOEJ, error) {
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

	return &McOEJ{
		tbc:            tbc,
		blockSize:      tbc.BlockSize(),
		blockBytes:     tbc.BlockSize() / 8,
		blockBytesMask: tbc.BlockSize()/8 - 1,
		blockUint64s:   tbc.BlockSize() / 64,
		ft:             &fullTweak{data: make([]byte, tbc.BlockSize()/8), mask: tbc.BlockSize()/8 - 1},
		blockBuf:       make([]byte, tbc.BlockSize()/8),
		fullKey:        make([]byte, tbc.BlockSize()/8+8),
	}, nil
}

type McOEJ struct {
	tbc            types.TweakableBlockCipher
	blockSize      int
	blockBytes     int
	blockBytesMask int
	blockUint64s   int
	ft             *fullTweak
	blockBuf       []byte
	fullKey        []byte
}

type fullTweak struct {
	data []byte
	buf  [24]byte
	pos  int
	mask int
}

func (f *fullTweak) Reset(nonce []byte) {
	copy(f.data, nonce)
	f.pos = 0
}

func (f *fullTweak) Xor(b []byte) {
	convert.XorBytes(f.data, f.data, b)
}

func (f *fullTweak) Next() (tweak []byte) {
	copy(f.buf[:], f.data[f.pos:f.pos+16])
	f.pos = (f.pos + 16) & f.mask
	return f.buf[:]
}

func (f *fullTweak) Remainder() []byte {
	return append(f.data, f.data...)[f.pos : f.pos+len(f.data)-16]
}

func (j *McOEJ) Lock(key, nonce, publicData, plaintext, dst []byte) []byte {
	if len(key) != j.blockBytes {
		panic("key length must equal block size")
	}
	if len(nonce) != j.blockBytes {
		panic("nonce length must equal block size")
	}
	// TODO: instead of hashing, we can run the encryption like normal on the public data, but not
	// include any of the cipher text in the output.  This way the tweak will be set according to
	// the public data.  I'm not sure how this will affect the proof.
	copy(j.fullKey, key)
	j.ft.Reset(nonce)
	j.processPublicData(j.fullKey, publicData)

	outputLen := len(plaintext)
	if outputLen&j.blockBytesMask != 0 {
		outputLen += j.blockBytes - (outputLen & j.blockBytesMask)
	}
	outputLen += j.blockBytes

	dstStart := len(dst)
	if dst == nil {
		dst = make([]byte, outputLen)
		copy(dst, plaintext)
	} else {
		dst = append(dst, make([]byte, outputLen)...)
		copy(dst[dstStart:], plaintext)
	}
	target := dst[dstStart:]

	for len(target) > j.blockBytes {
		tweak := j.ft.Next()
		j.ft.Xor(target[0:j.blockBytes])
		j.tbc.Encrypt(target[0:j.blockBytes], j.fullKey, tweak)
		j.ft.Xor(target[0:j.blockBytes])
		target = target[j.blockBytes:]
	}
	finalTweak := j.ft.Next()
	copy(target, j.ft.Remainder())
	if finalBlockLen := byte(len(plaintext) & j.blockBytesMask); finalBlockLen > 0 {
		target[j.blockBytes-16] = finalBlockLen
	} else {
		target[j.blockBytes-16] = byte(j.blockBytes)
	}
	j.tbc.Encrypt(target, j.fullKey, finalTweak)
	return dst
}

func (j *McOEJ) Unlock(key, nonce, publicData, ciphertext, dst []byte) ([]byte, error) {
	if len(key) != j.blockBytes {
		panic("key length must equal block size")
	}
	if len(nonce) != j.blockBytes {
		panic("nonce length must equal block size")
	}
	if len(ciphertext)&j.blockBytesMask != 0 {
		return nil, fmt.Errorf("ciphertext length isn't a multiple of the block length")
	}
	// TODO: instead of hashing, we can run the encryption like normal on the public data, but not
	// include any of the cipher text in the output.  This way the tweak will be set according to
	// the public data.  I'm not sure how this will affect the proof.
	copy(j.fullKey, key)
	j.ft.Reset(nonce)
	j.processPublicData(j.fullKey, publicData)

	dstStart := len(dst)
	if dst == nil {
		dst = make([]byte, len(ciphertext))
		copy(dst, ciphertext)
	} else {
		dst = append(dst, ciphertext...)
	}
	target := dst[dstStart:]

	for len(target) > j.blockBytes {
		tweak := j.ft.Next()
		j.ft.Xor(target[0:j.blockBytes])
		j.tbc.Decrypt(target[0:j.blockBytes], j.fullKey, tweak[:])
		j.ft.Xor(target[0:j.blockBytes])
		target = target[j.blockBytes:]
	}
	finalTweak := j.ft.Next()
	remainder := j.ft.Remainder()
	j.tbc.Decrypt(target, j.fullKey, finalTweak)
	if subtle.ConstantTimeCompare(target[0:len(remainder)], remainder) != 1 {
		return nil, fmt.Errorf("authentication failed")
	}
	finalBlockLen := int(target[len(remainder)])
	dst = dst[0 : len(dst)-2*j.blockBytes+finalBlockLen]
	return dst, nil
}

func (j *McOEJ) processPublicData(key, publicData []byte) {
	for len(publicData) >= j.blockBytes {
		copy(j.blockBuf, publicData)
		j.tbc.Encrypt(j.blockBuf, key, j.ft.Next())
		j.ft.Xor(j.blockBuf)
		publicData = publicData[j.blockBytes:]
	}
	copy(j.blockBuf, publicData)
	j.blockBuf[len(publicData)] = 1
	for i := len(publicData) + 1; i < len(j.blockBuf); i++ {
		j.blockBuf[i] = 0
	}
	j.tbc.Encrypt(j.blockBuf, key, j.ft.Next())
	j.ft.Xor(j.blockBuf)
}

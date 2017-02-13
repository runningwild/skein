package mcoej

import (
	"crypto/subtle"
	"fmt"

	"github.com/runningwild/skein/convert"
	"github.com/runningwild/skein/threefish/1024"
)

const (
	blockBytes1024 = 128
	blockMask1024  = blockBytes1024 - 1
	tweakBytes1024 = 24
)

func New1024() *McOEJ1024 {
	return &McOEJ1024{}
}

type McOEJ1024 struct {
	ft fullTweak1024
}

type fullTweak1024 struct {
	data [blockBytes1024]byte
	buf  [tweakBytes1024]byte
	rem  [2 * (blockBytes1024 - 16)]byte
	pos  int
}

func (f *fullTweak1024) Reset(nonce []byte) {
	copy(f.data[:], nonce)
	f.pos = 0
}

func (f *fullTweak1024) Xor(b []byte) {
	convert.XorBytes(f.data[:], f.data[:], b)
}

func (f *fullTweak1024) Next() (tweak []byte) {
	copy(f.buf[:], f.data[f.pos:f.pos+16])
	f.pos = (f.pos + 16) & blockMask1024
	return f.buf[:]
}

func (f *fullTweak1024) Remainder() []byte {
	copy(f.rem[:], f.data[:])
	copy(f.rem[blockBytes1024:], f.data[:])
	return f.rem[f.pos : f.pos+blockBytes1024-16]
}

func (j *McOEJ1024) Lock(key, nonce, publicData, plaintext, dst []byte) []byte {
	if len(key) != blockBytes1024 {
		panic("key length must equal block size")
	}
	if len(nonce) != blockBytes1024 {
		panic("nonce length must equal block size")
	}
	// TODO: instead of hashing, we can run the encryption like normal on the public data, but not
	// include any of the cipher text in the output.  This way the tweak will be set according to
	// the public data.  I'm not sure how this will affect the proof.
	var fullKey [blockBytes1024 + 8]byte
	copy(fullKey[:], key)
	j.ft.Reset(nonce)
	j.processPublicData(fullKey[:], publicData)

	outputLen := len(plaintext)
	if outputLen&blockMask1024 != 0 {
		outputLen += blockBytes1024 - (outputLen & blockMask1024)
	}
	outputLen += blockBytes1024

	dstStart := len(dst)
	if dst == nil {
		dst = make([]byte, outputLen)
		copy(dst, plaintext)
	} else {
		dst = append(dst, make([]byte, outputLen)...)
		copy(dst[dstStart:], plaintext)
	}
	target := dst[dstStart:]

	for len(target) > blockBytes1024 {
		tweak := j.ft.Next()
		j.ft.Xor(target[0:blockBytes1024])
		threefish.Encrypt(target[0:blockBytes1024], fullKey[:], tweak)
		j.ft.Xor(target[0:blockBytes1024])
		target = target[blockBytes1024:]
	}
	finalTweak := j.ft.Next()
	copy(target, j.ft.Remainder())
	if finalBlockLen := byte(len(plaintext) & blockMask1024); finalBlockLen > 0 {
		target[blockBytes1024-16] = finalBlockLen
	} else {
		target[blockBytes1024-16] = byte(blockBytes1024)
	}
	threefish.Encrypt(target, fullKey[:], finalTweak)
	return dst
}

func (j *McOEJ1024) Unlock(key, nonce, publicData, ciphertext, dst []byte) ([]byte, error) {
	if len(key) != blockBytes1024 {
		panic("key length must equal block size")
	}
	if len(nonce) != blockBytes1024 {
		panic("nonce length must equal block size")
	}
	if len(ciphertext)&blockMask1024 != 0 {
		return nil, fmt.Errorf("ciphertext length isn't a multiple of the block length")
	}
	// TODO: instead of hashing, we can run the encryption like normal on the public data, but not
	// include any of the cipher text in the output.  This way the tweak will be set according to
	// the public data.  I'm not sure how this will affect the proof.
	var fullKey [blockBytes1024 + 8]byte
	copy(fullKey[:], key)
	j.ft.Reset(nonce)
	j.processPublicData(fullKey[:], publicData)

	dstStart := len(dst)
	if dst == nil {
		dst = make([]byte, len(ciphertext))
		copy(dst, ciphertext)
	} else {
		dst = append(dst, ciphertext...)
	}
	target := dst[dstStart:]

	for len(target) > blockBytes1024 {
		tweak := j.ft.Next()
		j.ft.Xor(target[0:blockBytes1024])
		threefish.Decrypt(target[0:blockBytes1024], fullKey[:], tweak[:])
		j.ft.Xor(target[0:blockBytes1024])
		target = target[blockBytes1024:]
	}
	finalTweak := j.ft.Next()
	remainder := j.ft.Remainder()
	threefish.Decrypt(target, fullKey[:], finalTweak)
	if subtle.ConstantTimeCompare(target[0:len(remainder)], remainder) != 1 {
		return nil, fmt.Errorf("authentication failed")
	}
	finalBlockLen := int(target[len(remainder)])
	dst = dst[0 : len(dst)-2*blockBytes1024+finalBlockLen]
	return dst, nil
}

func (j *McOEJ1024) processPublicData(fullKey, publicData []byte) {
	var blockBuf [blockBytes1024]byte
	for len(publicData) >= blockBytes1024 {
		copy(blockBuf[:], publicData)
		threefish.Encrypt(blockBuf[:], fullKey, j.ft.Next())
		j.ft.Xor(blockBuf[:])
		publicData = publicData[blockBytes1024:]
	}
	copy(blockBuf[:], publicData)
	blockBuf[len(publicData)] = 1
	for i := len(publicData) + 1; i < len(blockBuf[:]); i++ {
		blockBuf[i] = 0
	}
	threefish.Encrypt(blockBuf[:], fullKey, j.ft.Next())
	j.ft.Xor(blockBuf[:])
}

package mcoej

import (
	"crypto/subtle"
	"fmt"

	"github.com/runningwild/skein/convert"
	"github.com/runningwild/skein/threefish/256"
)

const (
	blockBytes256 = 32
	blockMask256  = blockBytes256 - 1
	tweakBytes256 = 24
)

func New256() *McOEJ256 {
	return &McOEJ256{}
}

type McOEJ256 struct {
	ft fullTweak256
}

type fullTweak256 struct {
	data [blockBytes256]byte
	buf  [tweakBytes256]byte
	rem  [2 * (blockBytes256 - 16)]byte
	pos  int
}

func (f *fullTweak256) Reset(nonce []byte) {
	copy(f.data[:], nonce)
	f.pos = 0
}

func (f *fullTweak256) Xor(b []byte) {
	convert.XorBytes(f.data[:], f.data[:], b)
}

func (f *fullTweak256) Next() (tweak []byte) {
	copy(f.buf[:], f.data[f.pos:f.pos+16])
	f.pos = (f.pos + 16) & blockMask256
	return f.buf[:]
}

func (f *fullTweak256) Remainder() []byte {
	copy(f.rem[:], f.data[:])
	copy(f.rem[blockBytes256:], f.data[:])
	return f.rem[f.pos : f.pos+blockBytes256-16]
}

func (j *McOEJ256) Lock(key, nonce, publicData, plaintext, dst []byte) []byte {
	if len(key) != blockBytes256 {
		panic("key length must equal block size")
	}
	if len(nonce) != blockBytes256 {
		panic("nonce length must equal block size")
	}
	// TODO: instead of hashing, we can run the encryption like normal on the public data, but not
	// include any of the cipher text in the output.  This way the tweak will be set according to
	// the public data.  I'm not sure how this will affect the proof.
	var fullKey [blockBytes256 + 8]byte
	copy(fullKey[:], key)
	j.ft.Reset(nonce)
	j.processPublicData(fullKey[:], publicData)

	outputLen := len(plaintext)
	if outputLen&blockMask256 != 0 {
		outputLen += blockBytes256 - (outputLen & blockMask256)
	}
	outputLen += blockBytes256

	dstStart := len(dst)
	if dst == nil {
		dst = make([]byte, outputLen)
		copy(dst, plaintext)
	} else {
		dst = append(dst, make([]byte, outputLen)...)
		copy(dst[dstStart:], plaintext)
	}
	target := dst[dstStart:]

	for len(target) > blockBytes256 {
		tweak := j.ft.Next()
		j.ft.Xor(target[0:blockBytes256])
		threefish.Encrypt(target[0:blockBytes256], fullKey[:], tweak)
		j.ft.Xor(target[0:blockBytes256])
		target = target[blockBytes256:]
	}
	finalTweak := j.ft.Next()
	copy(target, j.ft.Remainder())
	if finalBlockLen := byte(len(plaintext) & blockMask256); finalBlockLen > 0 {
		target[blockBytes256-16] = finalBlockLen
	} else {
		target[blockBytes256-16] = byte(blockBytes256)
	}
	threefish.Encrypt(target, fullKey[:], finalTweak)
	return dst
}

func (j *McOEJ256) Unlock(key, nonce, publicData, ciphertext, dst []byte) ([]byte, error) {
	if len(key) != blockBytes256 {
		panic("key length must equal block size")
	}
	if len(nonce) != blockBytes256 {
		panic("nonce length must equal block size")
	}
	if len(ciphertext)&blockMask256 != 0 {
		return nil, fmt.Errorf("ciphertext length isn't a multiple of the block length")
	}
	// TODO: instead of hashing, we can run the encryption like normal on the public data, but not
	// include any of the cipher text in the output.  This way the tweak will be set according to
	// the public data.  I'm not sure how this will affect the proof.
	var fullKey [blockBytes256 + 8]byte
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

	for len(target) > blockBytes256 {
		tweak := j.ft.Next()
		j.ft.Xor(target[0:blockBytes256])
		threefish.Decrypt(target[0:blockBytes256], fullKey[:], tweak[:])
		j.ft.Xor(target[0:blockBytes256])
		target = target[blockBytes256:]
	}
	finalTweak := j.ft.Next()
	remainder := j.ft.Remainder()
	threefish.Decrypt(target, fullKey[:], finalTweak)
	if subtle.ConstantTimeCompare(target[0:len(remainder)], remainder) != 1 {
		return nil, fmt.Errorf("authentication failed")
	}
	finalBlockLen := int(target[len(remainder)])
	dst = dst[0 : len(dst)-2*blockBytes256+finalBlockLen]
	return dst, nil
}

func (j *McOEJ256) processPublicData(fullKey, publicData []byte) {
	var blockBuf [blockBytes256]byte
	for len(publicData) >= blockBytes256 {
		copy(blockBuf[:], publicData)
		threefish.Encrypt(blockBuf[:], fullKey, j.ft.Next())
		j.ft.Xor(blockBuf[:])
		publicData = publicData[blockBytes256:]
	}
	copy(blockBuf[:], publicData)
	blockBuf[len(publicData)] = 1
	for i := len(publicData) + 1; i < len(blockBuf[:]); i++ {
		blockBuf[i] = 0
	}
	threefish.Encrypt(blockBuf[:], fullKey, j.ft.Next())
	j.ft.Xor(blockBuf[:])
}

package mcoej

import (
	"crypto/subtle"
	"fmt"

	"github.com/runningwild/skein/convert"
	"github.com/runningwild/skein/threefish/512"
)

const (
	blockBytes512 = 64
	blockMask512  = blockBytes512 - 1
	tweakBytes512 = 24
)

func New512() *McOEJ512 {
	return &McOEJ512{}
}

type McOEJ512 struct {
	ft fullTweak512
}

type fullTweak512 struct {
	data [blockBytes512]byte
	buf  [tweakBytes512]byte
	rem  [2 * (blockBytes512 - 16)]byte
	pos  int
}

func (f *fullTweak512) Reset(nonce []byte) {
	copy(f.data[:], nonce)
	f.pos = 0
}

func (f *fullTweak512) Xor(b []byte) {
	convert.XorBytes(f.data[:], f.data[:], b)
}

func (f *fullTweak512) Next() (tweak []byte) {
	copy(f.buf[:], f.data[f.pos:f.pos+16])
	f.pos = (f.pos + 16) & blockMask512
	return f.buf[:]
}

func (f *fullTweak512) Remainder() []byte {
	copy(f.rem[:], f.data[:])
	copy(f.rem[blockBytes512:], f.data[:])
	return f.rem[f.pos : f.pos+blockBytes512-16]
}

func (j *McOEJ512) Lock(key, nonce, publicData, plaintext, dst []byte) []byte {
	if len(key) != blockBytes512 {
		panic("key length must equal block size")
	}
	if len(nonce) != blockBytes512 {
		panic("nonce length must equal block size")
	}
	// TODO: instead of hashing, we can run the encryption like normal on the public data, but not
	// include any of the cipher text in the output.  This way the tweak will be set according to
	// the public data.  I'm not sure how this will affect the proof.
	var fullKey [blockBytes512 + 8]byte
	copy(fullKey[:], key)
	j.ft.Reset(nonce)
	j.processPublicData(fullKey[:], publicData)

	outputLen := len(plaintext)
	if outputLen&blockMask512 != 0 {
		outputLen += blockBytes512 - (outputLen & blockMask512)
	}
	outputLen += blockBytes512

	dstStart := len(dst)
	if dst == nil {
		dst = make([]byte, outputLen)
		copy(dst, plaintext)
	} else {
		dst = append(dst, make([]byte, outputLen)...)
		copy(dst[dstStart:], plaintext)
	}
	target := dst[dstStart:]

	for len(target) > blockBytes512 {
		tweak := j.ft.Next()
		j.ft.Xor(target[0:blockBytes512])
		threefish.Encrypt(target[0:blockBytes512], fullKey[:], tweak)
		j.ft.Xor(target[0:blockBytes512])
		target = target[blockBytes512:]
	}
	finalTweak := j.ft.Next()
	copy(target, j.ft.Remainder())
	if finalBlockLen := byte(len(plaintext) & blockMask512); finalBlockLen > 0 {
		target[blockBytes512-16] = finalBlockLen
	} else {
		target[blockBytes512-16] = byte(blockBytes512)
	}
	threefish.Encrypt(target, fullKey[:], finalTweak)
	return dst
}

func (j *McOEJ512) Unlock(key, nonce, publicData, ciphertext, dst []byte) ([]byte, error) {
	if len(key) != blockBytes512 {
		panic("key length must equal block size")
	}
	if len(nonce) != blockBytes512 {
		panic("nonce length must equal block size")
	}
	if len(ciphertext)&blockMask512 != 0 {
		return nil, fmt.Errorf("ciphertext length isn't a multiple of the block length")
	}
	// TODO: instead of hashing, we can run the encryption like normal on the public data, but not
	// include any of the cipher text in the output.  This way the tweak will be set according to
	// the public data.  I'm not sure how this will affect the proof.
	var fullKey [blockBytes512 + 8]byte
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

	for len(target) > blockBytes512 {
		tweak := j.ft.Next()
		j.ft.Xor(target[0:blockBytes512])
		threefish.Decrypt(target[0:blockBytes512], fullKey[:], tweak[:])
		j.ft.Xor(target[0:blockBytes512])
		target = target[blockBytes512:]
	}
	finalTweak := j.ft.Next()
	remainder := j.ft.Remainder()
	threefish.Decrypt(target, fullKey[:], finalTweak)
	if subtle.ConstantTimeCompare(target[0:len(remainder)], remainder) != 1 {
		return nil, fmt.Errorf("authentication failed")
	}
	finalBlockLen := int(target[len(remainder)])
	dst = dst[0 : len(dst)-2*blockBytes512+finalBlockLen]
	return dst, nil
}

func (j *McOEJ512) processPublicData(fullKey, publicData []byte) {
	var blockBuf [blockBytes512]byte
	for len(publicData) >= blockBytes512 {
		copy(blockBuf[:], publicData)
		threefish.Encrypt(blockBuf[:], fullKey, j.ft.Next())
		j.ft.Xor(blockBuf[:])
		publicData = publicData[blockBytes512:]
	}
	copy(blockBuf[:], publicData)
	blockBuf[len(publicData)] = 1
	for i := len(publicData) + 1; i < len(blockBuf[:]); i++ {
		blockBuf[i] = 0
	}
	threefish.Encrypt(blockBuf[:], fullKey, j.ft.Next())
	j.ft.Xor(blockBuf[:])
}

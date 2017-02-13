package mcoej

import (
	"crypto/subtle"
	"fmt"

	"github.com/runningwild/skein/convert"
	"github.com/runningwild/skein/threefish/512"
)

func New512() *McOEJ512 {
	return &McOEJ512{}
}

type McOEJ512 struct {
	ft fullTweak512
}

type fullTweak512 struct {
	data [64]byte
	buf  [24]byte
	rem  [96]byte
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
	f.pos = (f.pos + 16) & 63
	return f.buf[:]
}

func (f *fullTweak512) Remainder() []byte {
	copy(f.rem[:], f.data[:])
	copy(f.rem[64:], f.data[:])
	return f.rem[f.pos : f.pos+48]
}

func (j *McOEJ512) Lock(key, nonce, publicData, plaintext, dst []byte) []byte {
	if len(key) != 64 {
		panic("key length must equal block size")
	}
	if len(nonce) != 64 {
		panic("nonce length must equal block size")
	}
	// TODO: instead of hashing, we can run the encryption like normal on the public data, but not
	// include any of the cipher text in the output.  This way the tweak will be set according to
	// the public data.  I'm not sure how this will affect the proof.
	var fullKey [64 + 8]byte
	copy(fullKey[:], key)
	j.ft.Reset(nonce)
	j.processPublicData(fullKey[:], publicData)

	outputLen := len(plaintext)
	if outputLen&63 != 0 {
		outputLen += 64 - (outputLen & 63)
	}
	outputLen += 64

	dstStart := len(dst)
	if dst == nil {
		dst = make([]byte, outputLen)
		copy(dst, plaintext)
	} else {
		dst = append(dst, make([]byte, outputLen)...)
		copy(dst[dstStart:], plaintext)
	}
	target := dst[dstStart:]

	for len(target) > 64 {
		tweak := j.ft.Next()
		j.ft.Xor(target[0:64])
		threefish.Encrypt(target[0:64], fullKey[:], tweak)
		j.ft.Xor(target[0:64])
		target = target[64:]
	}
	finalTweak := j.ft.Next()
	copy(target, j.ft.Remainder())
	if finalBlockLen := byte(len(plaintext) & 63); finalBlockLen > 0 {
		target[64-16] = finalBlockLen
	} else {
		target[64-16] = byte(64)
	}
	threefish.Encrypt(target, fullKey[:], finalTweak)
	return dst
}

func (j *McOEJ512) Unlock(key, nonce, publicData, ciphertext, dst []byte) ([]byte, error) {
	if len(key) != 64 {
		panic("key length must equal block size")
	}
	if len(nonce) != 64 {
		panic("nonce length must equal block size")
	}
	if len(ciphertext)&63 != 0 {
		return nil, fmt.Errorf("ciphertext length isn't a multiple of the block length")
	}
	// TODO: instead of hashing, we can run the encryption like normal on the public data, but not
	// include any of the cipher text in the output.  This way the tweak will be set according to
	// the public data.  I'm not sure how this will affect the proof.
	var fullKey [64 + 8]byte
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

	for len(target) > 64 {
		tweak := j.ft.Next()
		j.ft.Xor(target[0:64])
		threefish.Decrypt(target[0:64], fullKey[:], tweak[:])
		j.ft.Xor(target[0:64])
		target = target[64:]
	}
	finalTweak := j.ft.Next()
	remainder := j.ft.Remainder()
	threefish.Decrypt(target, fullKey[:], finalTweak)
	if subtle.ConstantTimeCompare(target[0:len(remainder)], remainder) != 1 {
		return nil, fmt.Errorf("authentication failed")
	}
	finalBlockLen := int(target[len(remainder)])
	dst = dst[0 : len(dst)-2*64+finalBlockLen]
	return dst, nil
}

func (j *McOEJ512) processPublicData(fullKey, publicData []byte) {
	var blockBuf [64]byte
	for len(publicData) >= 64 {
		copy(blockBuf[:], publicData)
		threefish.Encrypt(blockBuf[:], fullKey, j.ft.Next())
		j.ft.Xor(blockBuf[:])
		publicData = publicData[64:]
	}
	copy(blockBuf[:], publicData)
	blockBuf[len(publicData)] = 1
	for i := len(publicData) + 1; i < len(blockBuf[:]); i++ {
		blockBuf[i] = 0
	}
	threefish.Encrypt(blockBuf[:], fullKey, j.ft.Next())
	j.ft.Xor(blockBuf[:])
}

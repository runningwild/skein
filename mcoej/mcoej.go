package mcoej

import (
	"crypto/subtle"
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

type fullTweak struct {
	data []byte
	buf  [24]byte
	pos  int
}

func (f *fullTweak) Xor(b []byte) {
	convert.XorBytes(f.data, f.data, b)
}

func (f *fullTweak) Next() (tweak []byte) {
	copy(f.buf[:], f.data[f.pos:f.pos+16])
	f.pos = (f.pos + 16) % len(f.data)
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
	fullKey := make([]byte, len(key)+8)
	copy(fullKey, key)
	ft := &fullTweak{data: make([]byte, j.blockBytes)}
	copy(ft.data, nonce)
	buf := make([]byte, j.blockBytes)
	for len(publicData) >= j.blockBytes {
		copy(buf, publicData)
		j.enc(buf, fullKey, ft.Next())
		ft.Xor(buf)
		publicData = publicData[j.blockBytes:]
	}
	copy(buf, publicData)
	buf[len(publicData)] = 1
	for i := len(publicData) + 1; i < len(buf); i++ {
		buf[i] = 0
	}
	j.enc(buf, fullKey, ft.Next())
	ft.Xor(buf)

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
		dst = append(dst, plaintext...)
	}
	target := dst[dstStart:]

	for len(target) > j.blockBytes {
		tweak := ft.Next()
		ft.Xor(target[0:j.blockBytes])
		j.enc(target[0:j.blockBytes], fullKey, tweak)
		ft.Xor(target[0:j.blockBytes])
		target = target[j.blockBytes:]
	}
	finalTweak := ft.Next()
	copy(target, ft.Remainder())
	target[j.blockBytes-16] = byte(len(plaintext) & j.blockBytesMask)
	j.enc(target, fullKey, finalTweak)
	return dst
}

func (j *McOEJ) Unlock(key, nonce, publicData, ciphertext, dst []byte) ([]byte, error) {
	if len(key) != j.blockBytes {
		panic("key length must equal block size")
	}
	if len(nonce) != j.blockBytes {
		panic("nonce length must equal block size")
	}
	// TODO: instead of hashing, we can run the encryption like normal on the public data, but not
	// include any of the cipher text in the output.  This way the tweak will be set according to
	// the public data.  I'm not sure how this will affect the proof.
	fullKey := make([]byte, len(key)+8)
	copy(fullKey, key)
	ft := &fullTweak{data: make([]byte, j.blockBytes)}
	copy(ft.data, nonce)
	buf := make([]byte, j.blockBytes)
	for len(publicData) >= j.blockBytes {
		copy(buf, publicData)
		j.enc(buf, fullKey, ft.Next())
		ft.Xor(buf)
		publicData = publicData[j.blockBytes:]
	}
	copy(buf, publicData)
	buf[len(publicData)] = 1
	for i := len(publicData) + 1; i < len(buf); i++ {
		buf[i] = 0
	}
	j.enc(buf, fullKey, ft.Next())
	ft.Xor(buf)

	dstStart := len(dst)
	if dst == nil {
		dst = make([]byte, len(ciphertext))
		copy(dst, ciphertext)
	} else {
		dst = append(dst, ciphertext...)
	}
	target := dst[dstStart:]

	for len(target) > j.blockBytes {
		tweak := ft.Next()
		ft.Xor(target[0:j.blockBytes])
		j.dec(target[0:j.blockBytes], fullKey, tweak[:])
		ft.Xor(target[0:j.blockBytes])
		target = target[j.blockBytes:]
	}
	finalTweak := ft.Next()
	remainder := ft.Remainder()
	j.dec(target, fullKey, finalTweak)
	if subtle.ConstantTimeCompare(target[0:len(remainder)], remainder) != 1 {
		return nil, fmt.Errorf("authentication failed")
	}
	finalBlockLen := int(target[len(remainder)])
	dst = dst[0 : len(dst)-2*j.blockBytes+finalBlockLen]
	return dst, nil
}

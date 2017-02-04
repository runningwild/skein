// package threefish is a convenience wrapper around the three specifications in the subfolders 256,
// 512, and 1024.
package threefish

import (
	"crypto/cipher"
	"fmt"

	tf1024 "github.com/runningwild/skein/threefish/1024"
	tf256 "github.com/runningwild/skein/threefish/256"
	tf512 "github.com/runningwild/skein/threefish/512"
)

// NewCipher returns a cipher.Block object using the either threefish-256, threefish-512, or
// threefish-1024, depending on the key size.
func NewCipher(key []byte) (cipher.Block, error) {
	if len(key) == 32 {
		var a [32]byte
		copy(a[:], key)
		return tf256.MakeCipher(a), nil
	}
	if len(key) == 64 {
		var a [64]byte
		copy(a[:], key)
		return tf512.MakeCipher(a), nil
	}
	if len(key) == 128 {
		var a [128]byte
		copy(a[:], key)
		return tf1024.MakeCipher(a), nil
	}
	return nil, fmt.Errorf("invalid key length %d, must be one of 256, 512, or 1024", len(key))
}

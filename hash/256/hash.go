// Package hash provides functions for computing skein hashes using 256-bit skein.
package hash

import (
	"fmt"
	"hash"

	"github.com/runningwild/skein/hash/hasher"
	"github.com/runningwild/skein/threefish/256"
	"github.com/runningwild/skein/ubi"
)

var (
	u *ubi.UBI
)

func init() {
	var err error
	if u, err = ubi.New(threefish.TweakableBlockCipher{}); err != nil {
		panic(fmt.Sprintf("failed to create ubi object: %v", err))
	}
}

// NewHash256 returns a hash.Hash object that computes N-bit hashes using 256-bit skein.
func NewHash256(N int) hash.Hash {
	return hasher.NewHasher(u, N)
}

// NewMAC256 returns a hash.Hash object that computes N-bit MAC using key and 256-bit skein.
func NewMAC256(N int, key []byte) hash.Hash {
	return hasher.NewMACer(u, key, N)
}

// Hash256 returns the N-bit hash of data using 256-bit skein.
func Hash256(N int, data []byte) []byte {
	return hasher.Hash(u, data, 0, uint64(N))
}

// TreeHash256 returns the N-bit tree hash of data using 256-bit skein and tree parameters Yl, Yf, and Ym.
func TreeHash256(N int, Yl, Yf, Ym byte, data []byte) []byte {
	return hasher.TreeHash(u, data, 0, uint64(N), Yl, Yf, Ym)
}

// MAC256 returns the N-bit MAC of data using key 256-bit skein.
func MAC256(N int, key, data []byte) []byte {
	return hasher.MAC(u, key, data, 0, uint64(N))
}

// Hash256Bits returns the N-bit hash of data using 256-bit skein.  lbb specifies the number of bits
// that should be used in the last byte of data.
func Hash256Bits(N int, data []byte, lbb int) []byte {
	return hasher.Hash(u, data, lbb, uint64(N))
}

// MAC256Bits returns the N-bit MAC of data using key 256-bit skein.  lbb specifies the number of bits
// that should be used in the last byte of data.
func MAC256Bits(N int, key, data []byte, lbb int) []byte {
	return hasher.MAC(u, key, data, lbb, uint64(N))
}

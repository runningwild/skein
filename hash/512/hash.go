// Package hash provides functions for computing skein hashes using 512-bit skein.
package hash

import (
	"fmt"
	"hash"

	"github.com/runningwild/skein/hash/hasher"
	"github.com/runningwild/skein/threefish/512"
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

// NewHash512 returns a hash.Hash object that computes N-bit hashes using 512-bit skein.
func NewHash512(N int) hash.Hash {
	return hasher.NewHasher(u, N)
}

// NewTreeHash512 reurns a hash.Hash object that computes a 512-bit Skein tree hash using the
// specified parameters.
func NewTreeHash512(N int, Yl, Yf, Ym byte) hash.Hash {
	return hasher.NewTreeHasher(u, N, Yl, Yf, Ym)
}

// NewMAC512 returns a hash.Hash object that computes N-bit MAC using key and 512-bit skein.
func NewMAC512(N int, key []byte) hash.Hash {
	return hasher.NewMACer(u, key, N)
}

// Hash512 returns the N-bit hash of data using 512-bit skein.
func Hash512(N int, data []byte) []byte {
	return hasher.Hash(u, data, 0, uint64(N))
}

// TreeHash512 returns the N-bit tree hash of data using 512-bit skein and tree parameters Yl, Yf, and Ym.
func TreeHash512(N int, Yl, Yf, Ym byte, data []byte) []byte {
	return hasher.TreeHash(u, data, 0, uint64(N), Yl, Yf, Ym)
}

// MAC512 returns the N-bit MAC of data using key 512-bit skein.
func MAC512(N int, key, data []byte) []byte {
	return hasher.MAC(u, key, data, 0, uint64(N))
}

// Hash512Bits returns the N-bit hash of data using 512-bit skein.  lbb specifies the number of bits
// that should be used in the last byte of data.
func Hash512Bits(N int, data []byte, lbb int) []byte {
	return hasher.Hash(u, data, lbb, uint64(N))
}

// MAC512Bits returns the N-bit MAC of data using key 512-bit skein.  lbb specifies the number of bits
// that should be used in the last byte of data.
func MAC512Bits(N int, key, data []byte, lbb int) []byte {
	return hasher.MAC(u, key, data, lbb, uint64(N))
}

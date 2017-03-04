// Package hash provides functions for computing skein hashes using 512-bit skein.
package hash

import (
	"fmt"
	"hash"

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
	return u.NewHasher(N)
}

// NewMAC512 returns a hash.Hash object that computes N-bit MAC using key and 512-bit skein.
func NewMAC512(N int, key []byte) hash.Hash {
	return u.NewMACer(key, N)
}

// Hash512 returns the N-bit hash of data using 512-bit skein.
func Hash512(N int, data []byte) []byte {
	return u.Hash(data, len(data)*8, uint64(N))
}

// MAC512 returns the N-bit MAC of data using key 512-bit skein.
func MAC512(N int, key, data []byte) []byte {
	return u.MAC(key, data, len(data)*8, uint64(N))
}

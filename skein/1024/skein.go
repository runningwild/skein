// Package skein provides functions for computing skein hashes using 1024-bit skein.
package skein

import (
	"fmt"
	"hash"

	"github.com/runningwild/skein/threefish/1024"
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

// NewHash1024 returns a hash.Hash object that computes N-bit hashes using 1024-bit skein.
func NewHash1024(N int) hash.Hash {
	return u.NewHasher(N)
}

// NewMAC1024 returns a hash.Hash object that computes N-bit MAC using key and 1024-bit skein.
func NewMAC1024(N int, key []byte) hash.Hash {
	return u.NewMACer(key, N)
}

// Hash1024 returns the N-bit hash of data using 1024-bit skein.
func Hash1024(N int, data []byte) []byte {
	return u.Hash(data, len(data)*8, uint64(N))
}

// MAC1024 returns the N-bit MAC of data using key 1024-bit skein.
func MAC1024(N int, key, data []byte) []byte {
	return u.MAC(key, data, len(data)*8, uint64(N))
}

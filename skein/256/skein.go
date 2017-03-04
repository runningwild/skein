// Package skein provides functions for computing skein hashes using 256-bit skein.
package skein

import (
	"fmt"
	"hash"

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
	return u.NewHasher(N)
}

// NewMAC256 returns a hash.Hash object that computes N-bit MAC using key and 256-bit skein.
func NewMAC256(N int, key []byte) hash.Hash {
	return u.NewMACer(key, N)
}

// Hash256 returns the N-bit hash of data using 256-bit skein.
func Hash256(N int, data []byte) []byte {
	return u.Hash(data, len(data)*8, uint64(N))
}

// MAC256 returns the N-bit MAC of data using key 256-bit skein.
func MAC256(N int, key, data []byte) []byte {
	return u.MAC(key, data, len(data)*8, uint64(N))
}

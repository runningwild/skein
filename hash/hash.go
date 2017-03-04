// Package hash is a convenience package that includes all three version of skein.
package hash

import (
	"hash"

	hash1024 "github.com/runningwild/skein/hash/1024"
	hash256 "github.com/runningwild/skein/hash/256"
	hash512 "github.com/runningwild/skein/hash/512"
)

// NewHash256 returns a hash.Hash object that computes N-bit hashes using 256-bit skein.
func NewHash256(N int) hash.Hash {
	return hash256.NewHash256(N)
}

// NewMAC256 returns a hash.Hash object that computes N-bit MAC using key and 256-bit skein.
func NewMAC256(N int, key []byte) hash.Hash {
	return hash256.NewMAC256(N, key)
}

// Hash256 returns the N-bit hash of data using 256-bit skein.
func Hash256(N int, data []byte) []byte {
	return hash256.Hash256(N, data)
}

// MAC256 returns the N-bit MAC of data using key 256-bit skein.
func MAC256(N int, key, data []byte) []byte {
	return hash256.MAC256(N, key, data)
}

// NewHash512 returns a hash.Hash object that computes N-bit hashes using 512-bit skein.
func NewHash512(N int) hash.Hash {
	return hash512.NewHash512(N)
}

// NewMAC512 returns a hash.Hash object that computes N-bit MAC using key and 512-bit skein.
func NewMAC512(N int, key []byte) hash.Hash {
	return hash512.NewMAC512(N, key)
}

// Hash512 returns the N-bit hash of data using 512-bit skein.
func Hash512(N int, data []byte) []byte {
	return hash512.Hash512(N, data)
}

// MAC512 returns the N-bit MAC of data using key 512-bit skein.
func MAC512(N int, key, data []byte) []byte {
	return hash512.MAC512(N, key, data)
}

// NewHash1024 returns a hash.Hash object that computes N-bit hashes using 1024-bit skein.
func NewHash1024(N int) hash.Hash {
	return hash1024.NewHash1024(N)
}

// NewMAC1024 returns a hash.Hash object that computes N-bit MAC using key and 1024-bit skein.
func NewMAC1024(N int, key []byte) hash.Hash {
	return hash1024.NewMAC1024(N, key)
}

// Hash1024 returns the N-bit hash of data using 1024-bit skein.
func Hash1024(N int, data []byte) []byte {
	return hash1024.Hash1024(N, data)
}

// MAC1024 returns the N-bit MAC of data using key 1024-bit skein.
func MAC1024(N int, key, data []byte) []byte {
	return hash1024.MAC1024(N, key, data)
}

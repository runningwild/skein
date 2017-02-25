// Package types defines types common to multiple other packages.
package types

type TweakableBlockCipher interface {
	// Encrypt encrypts data using key and tweak.  key and tweak must be of the appropriate size or this
	// function will panic.
	Encrypt(data []byte, key []byte, tweak []byte)

	// Decrypt decrypts data using key and tweak.  key and tweak must be of the appropriate size or this
	// function will panic.
	Decrypt(data []byte, key []byte, tweak []byte)

	// BlockSize returns block size of this cipher in bits.
	BlockSize() int

	// TweakSize returns tweak size of this cipher in bits.
	TweakSize() int
}

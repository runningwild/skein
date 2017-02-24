package skein

import (
	skein1024 "github.com/runningwild/skein/hash/1024"
	skein256 "github.com/runningwild/skein/hash/256"
	skein512 "github.com/runningwild/skein/hash/512"
)

var (
	Hash256  = skein256.Hash256
	MAC256   = skein256.MAC256
	Hash512  = skein512.Hash512
	MAC512   = skein512.MAC512
	Hash1024 = skein1024.Hash1024
	MAC1024  = skein1024.MAC1024
)

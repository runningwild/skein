package skein

import (
	"fmt"

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

func Hash256(data []byte, N int) []byte {
	return u.Hash(data, len(data)*8, uint64(N))
}

func MAC256(key []byte, data []byte, N int) []byte {
	return u.MAC(key, data, len(data)*8, uint64(N))
}
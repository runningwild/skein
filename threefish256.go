// +build amd64

package skein

type block256 struct {
	state [4]uint64
	key   [5]uint64 // The user-defined key is the first 4 values, we add the 5th.
	tweak [3]uint64 // The user-defined tweak is the first 2 values, we add the 3rd.
}

func (b *block256) encrypt() {
	encrypt256(&b.state, &b.key, &b.tweak)
}

func (b *block256) decrypt() {
	decrypt256(&b.state, &b.key, &b.tweak)
}

func encrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64)
func decrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64)

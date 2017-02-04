// +build !arm

package threefish

func encrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64) {
	encrypt256Simple(state, key, tweak)
}
func decrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64) {
	decrypt256Simple(state, key, tweak)
}

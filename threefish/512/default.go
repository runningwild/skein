package threefish

func encrypt512(state *[8]uint64, key *[9]uint64, tweak *[3]uint64) {
	encrypt512Simple(state, key, tweak)
}
func decrypt512(state *[8]uint64, key *[9]uint64, tweak *[3]uint64) {
	decrypt512Simple(state, key, tweak)
}
